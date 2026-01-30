package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/redis/go-redis/v9"
)

func (s *Service) VerifyAccount(ctx context.Context, data *dto.VerifyAccountRequest) error {
	key := s.generateVerificationKey(data.Token)

	userId, err := s.redis.GetDel(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if strings.TrimSpace(userId) == "" {
		return apperror.New(apperror.Conflict, "invalid token")
	}

	user, err := s.repository.GetById(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil || user.Email != data.Email || user.IsBanned {
		return apperror.New(apperror.Conflict, "invalid token")
	}

	user.IsVerified = true
	return s.repository.Update(ctx, user)
}
