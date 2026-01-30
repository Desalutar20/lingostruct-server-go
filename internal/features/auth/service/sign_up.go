package service

import (
	"context"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/Desalutar20/lingostruct-server-go/pkg/random"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SignUp(ctx context.Context, data *dto.SignUpRequest) error {
	existingUser, err := s.repository.GetByEmail(ctx, data.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return apperror.New(apperror.Conflict, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		FirstName:      data.FirstName,
		LastName:       data.LastName,
		Email:          data.Email,
		HashedPassword: string(hashedPassword),
	}

	id, err := s.repository.Create(ctx, &user)
	if err != nil {
		return err
	}

	secureToken, err := random.GenerateSecureToken(46)
	if err != nil {
		return err
	}

	if err := s.redis.SetEx(ctx, s.generateVerificationKey(secureToken), id, time.Minute*time.Duration(s.config.AccountVerificationTTLMinutes)).Err(); err != nil {
		return err
	}

	if err := s.sendVerificationEmail(secureToken, data.Email); err != nil {
		return err
	}

	return nil
}
