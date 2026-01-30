package service

import (
	"context"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	userDto "github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"
	"golang.org/x/crypto/bcrypt"

	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
)

func (s *Service) SignIn(ctx context.Context, data *dto.SignInRequest, w http.ResponseWriter) (*dto.UserWithSessionId, error) {
	user, err := s.repository.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	if user == nil || user.IsBanned || !user.IsVerified {
		return nil, apperror.New(apperror.Conflict, "invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(data.Password)); err != nil {
		return nil, apperror.New(apperror.Conflict, "invalid email or password")
	}

	sessionId, err := s.generateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, s.generateSessionCookieOptions(sessionId))
	return &dto.UserWithSessionId{
		SessionId: sessionId,
		UserResponse: userDto.UserResponse{
			ID:        user.ID,
			DeletedAt: user.DeletedAt,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			AvatarId:  user.AvatarId,
			AvatarUrl: user.AvatarUrl,
		},
	}, nil
}
