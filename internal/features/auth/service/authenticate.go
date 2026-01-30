package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/redis/go-redis/v9"
)

func (s *Service) Authenticate(ctx context.Context, sessionId string, w http.ResponseWriter) (*dto.UserResponse, error) {
	userId, err := s.redis.GetEx(ctx, s.generateSessionKey(sessionId), time.Minute*time.Duration(s.config.SessionTTLMinutes)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if strings.TrimSpace(userId) == "" {
		return nil, apperror.New(apperror.Unauthorized, "unauthorized")
	}

	user, err := s.repository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsBanned || !user.IsVerified {
		return nil, apperror.New(apperror.Conflict, "unauthorized")
	}

	http.SetCookie(w, s.generateSessionCookieOptions(sessionId))
	return &dto.UserResponse{
		ID:        user.ID,
		DeletedAt: user.DeletedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		AvatarId:  user.AvatarId,
		AvatarUrl: user.AvatarUrl,
	}, nil
}
