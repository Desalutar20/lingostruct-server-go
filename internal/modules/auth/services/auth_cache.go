package authservices

import (
	"context"
	"errors"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokenType int

const (
	accountVerification tokenType = iota
	resetPassword
)

type RedisKeyGenerator interface {
	AccountVerificationKey(token valueobjects.NonEmptyString) string
	SessionKey(sessionId uuid.UUID) string
	SessionsKey(userId userdomain.UserId) string
	PasswordResetKey(token valueobjects.NonEmptyString) string
}

type AuthCache struct {
	redis        *redis.Client
	appConfig    *config.ApplicationConfig
	keyGenerator RedisKeyGenerator
}

func (c *AuthCache) StoreUserIdByVerificationToken(ctx context.Context, userId userdomain.UserId, token valueobjects.NonEmptyString) error {
	return c.storeUserId(ctx, accountVerification, userId, token)
}

func (c *AuthCache) StoreUserIdByResetPasswordToken(ctx context.Context, userId userdomain.UserId, token valueobjects.NonEmptyString) error {
	return c.storeUserId(ctx, resetPassword, userId, token)
}

func (c *AuthCache) GetAndDeleteUserIdByVerificationToken(ctx context.Context, token valueobjects.NonEmptyString) (*userdomain.UserId, error) {
	return c.getUserId(ctx, accountVerification, token)
}

func (c *AuthCache) GetAndDeleteUserIdByResetPasswordToken(ctx context.Context, token valueobjects.NonEmptyString) (*userdomain.UserId, error) {
	return c.getUserId(ctx, resetPassword, token)
}

func (c *AuthCache) storeUserId(ctx context.Context, tokenType tokenType, userId userdomain.UserId, token valueobjects.NonEmptyString) error {
	var key string

	switch tokenType {
	case accountVerification:
		key = c.keyGenerator.AccountVerificationKey(token)
	case resetPassword:
		key = c.keyGenerator.PasswordResetKey(token)
	}

	return c.redis.Set(ctx, key, userId, c.appConfig.ResetPasswordTtl).Err()
}

func (c *AuthCache) getUserId(ctx context.Context, tokenType tokenType, token valueobjects.NonEmptyString) (*userdomain.UserId, error) {
	var key string

	switch tokenType {
	case accountVerification:
		key = c.keyGenerator.AccountVerificationKey(token)
	case resetPassword:
		key = c.keyGenerator.PasswordResetKey(token)
	}

	value, err := c.redis.GetDel(ctx, key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, domain.InternalError("Invalid user id stored in Redis")
	}

	return new(userdomain.UserId(parsed)), nil
}

func NewAuthCache(redis *redis.Client, appConfig *config.ApplicationConfig, keyGenerator RedisKeyGenerator) *AuthCache {
	return &AuthCache{
		redis,
		appConfig,
		keyGenerator,
	}
}
