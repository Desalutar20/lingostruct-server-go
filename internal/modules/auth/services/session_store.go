package authservices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authtypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/types"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStore struct {
	redis        *redis.Client
	keyGenerator RedisKeyGenerator
	appConfig    *config.ApplicationConfig
	redisConfig  *config.RedisConfig
}

func (s *SessionStore) Save(ctx context.Context, session authtypes.Session) (uuid.UUID, error) {
	sessionId := uuid.New()

	sessionKey := s.keyGenerator.SessionKey(sessionId)
	sessionsKey := s.keyGenerator.SessionsKey(session.UserId)

	data, err := json.Marshal(session)
	if err != nil {
		return uuid.Nil, domain.InternalError(
			fmt.Sprintf("SessionStore.Save: failed to marshal session: %v", err),
		)
	}

	ttl := s.appConfig.SessionTtl
	expiresAt := time.Now().Add(ttl)
	score := float64(expiresAt.UnixMilli())

	pipe := s.redis.TxPipeline()

	pipe.Set(ctx, sessionKey, data, ttl)
	pipe.ZAdd(ctx, sessionsKey, redis.Z{Score: score, Member: sessionId.String()})

	if _, err := pipe.Exec(ctx); err != nil {
		return uuid.Nil, domain.InternalError(
			fmt.Sprintf("SessionStore.Save: failed to execute Redis pipeline: %v", err),
		)
	}

	return sessionId, nil
}

func (s *SessionStore) Get(ctx context.Context, sessionId uuid.UUID) (*authtypes.Session, error) {
	sessionKey := s.keyGenerator.SessionKey(sessionId)

	data, err := s.redis.Get(ctx, sessionKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, domain.InternalError(
			fmt.Sprintf("SessionStore.Get: failed to get session: %v", err),
		)
	}

	var session authtypes.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("SessionStore.Get: failed to get session: %v", err),
		)
	}

	return &session, nil
}

func (s *SessionStore) GetSessionIds(ctx context.Context, userId userdomain.UserId) ([]uuid.UUID, error) {
	sessionsKey := s.keyGenerator.SessionsKey(userId)
	now := time.Now().UnixMilli()

	values, err := s.redis.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:     sessionsKey,
		ByScore: true,
		Start:   strconv.FormatInt(now, 10),
		Stop:    "+inf",
	}).Result()
	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("SessionStore.GetSessionIds: failed to retrieve session IDs: %v", err),
		)
	}

	ids := make([]uuid.UUID, 0, len(values))

	for _, v := range values {
		member, ok := v.Member.(string)
		if !ok {
			return nil, domain.InternalError(
				fmt.Sprintf("SessionStore.GetSessionIds: unexpected session ID type: %T", v.Member),
			)
		}

		id, err := uuid.Parse(member)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (s *SessionStore) Remove(ctx context.Context, userId userdomain.UserId, sessionId uuid.UUID) error {
	sessionKey := s.keyGenerator.SessionKey(sessionId)
	sessionsKey := s.keyGenerator.SessionsKey(userId)

	pipe := s.redis.TxPipeline()

	pipe.Del(ctx, sessionKey)
	pipe.ZRem(ctx, sessionsKey, sessionId.String())

	if _, err := pipe.Exec(ctx); err != nil {
		return domain.InternalError(
			fmt.Sprintf("SessionStore.Remove: failed to execute Redis pipeline: %v", err),
		)
	}

	return nil
}

func (s *SessionStore) RemoveAll(ctx context.Context, userId userdomain.UserId) error {
	sessionsKey := s.keyGenerator.SessionsKey(userId)

	sessionIds, err := s.GetSessionIds(ctx, userId)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(sessionIds)+1)

	for _, id := range sessionIds {
		keys = append(keys, s.keyGenerator.SessionKey(id))
	}

	keys = append(keys, sessionsKey)

	if err := s.redis.Del(ctx, keys...).Err(); err != nil {
		return domain.InternalError(fmt.Sprintf("SessionStore.RemoveAll: failed to get session IDs: %v", err))
	}

	return nil
}

func (s *SessionStore) RemoveExpired(ctx context.Context) error {
	sessionsKeyPattern := fmt.Sprintf(
		"%s%s*",
		s.redisConfig.KeyPrefix,
		s.keyGenerator.SessionsPrefix(),
	)

	sessionKeyPrefix := fmt.Sprintf(
		"%s%s",
		s.redisConfig.KeyPrefix,
		s.keyGenerator.SessionPrefix(),
	)

	now := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var cursor uint64

	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, sessionsKeyPattern, 100).Result()
		if err != nil {
			return domain.InternalError(fmt.Sprintf("SessionStore.RemoveExpired: failed to scan sessions: %v", err))
		}

		for _, sessionsKey := range keys {
			sessions, err := s.redis.ZRangeArgs(ctx, redis.ZRangeArgs{
				Key:     sessionsKey,
				Start:   "-inf",
				Stop:    now,
				ByScore: true,
			}).Result()
			if err != nil {
				return domain.InternalError(fmt.Sprintf(
					"SessionStore.RemoveExpired: failed to retrieve expired sessions from %s: %v",
					sessionsKey, err,
				))
			}

			if len(sessions) == 0 {
				continue
			}

			keysToDelete := make([]string, 0, len(sessions))

			for _, session := range sessions {
				keysToDelete = append(
					keysToDelete,
					sessionKeyPrefix+session,
				)
			}

			pipe := s.redis.Pipeline()

			pipe.Del(ctx, keysToDelete...)
			pipe.ZRem(ctx, sessionsKey, membersToInterfaces(sessions)...)

			if _, err := pipe.Exec(ctx); err != nil {
				return domain.InternalError(fmt.Sprintf(
					"SessionStore.RemoveExpired: failed to delete expired sessions from %s: %v",
					sessionsKey, err,
				))
			}

			length, err := s.redis.ZCard(ctx, sessionsKey).Result()
			if err != nil {
				return domain.InternalError(fmt.Sprintf(
					"SessionStore.RemoveExpired: failed to count sessions in %s: %v",
					sessionsKey, err,
				))
			}

			if length == 0 {
				if err := s.redis.Del(ctx, sessionsKey).Err(); err != nil {
					return domain.InternalError(fmt.Sprintf(
						"SessionStore.RemoveExpired: failed to delete empty sessions set %s: %v",
						sessionsKey, err,
					))
				}
			}
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}

func membersToInterfaces(members []string) []interface{} {
	result := make([]interface{}, len(members))

	for i, member := range members {
		result[i] = member
	}

	return result
}

func NewSessionStore(redis *redis.Client, keyGenerator RedisKeyGenerator, appConfig *config.ApplicationConfig, redisConfig *config.RedisConfig) *SessionStore {
	return &SessionStore{
		redis,
		keyGenerator,
		appConfig,
		redisConfig,
	}
}
