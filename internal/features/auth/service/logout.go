package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func (s *Service) Logout(ctx context.Context, sessionId string) error {
	if err := s.redis.Del(ctx, s.generateSessionKey(sessionId)).Err(); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}
