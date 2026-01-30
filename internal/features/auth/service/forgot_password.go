package service

import (
	"context"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/pkg/random"
)

func (s *Service) ForgotPassword(ctx context.Context, data *dto.ForgotPasswordRequest) error {
	user, err := s.repository.GetByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	if user == nil || user.IsBanned || !user.IsVerified {
		return nil
	}

	secureToken, err := random.GenerateSecureToken(46)
	if err != nil {
		return err
	}

	if err := s.redis.SetEx(ctx, s.generateResetPasswordkey(secureToken), user.ID, time.Minute*time.Duration(s.config.ResetPasswordTTLMinutes)).Err(); err != nil {
		return err
	}

	if err := s.sendResetPasswordEmail(secureToken, data.Email); err != nil {
		return err
	}

	return nil
}
