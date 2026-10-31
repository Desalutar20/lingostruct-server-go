package services

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

type RedisKeyGenerator struct {
	keyPrefix string
}

func (g *RedisKeyGenerator) AccountVerificationPrefix() string {
	return "account-verification:"
}

func (g *RedisKeyGenerator) PasswordResetPrefix() string {
	return "password-reset:"
}

func (g *RedisKeyGenerator) SessionPrefix() string {
	return "session:"
}

func (g *RedisKeyGenerator) SessionsPrefix() string {
	return "sessions:"
}

func (g *RedisKeyGenerator) AccountVerificationKey(token valueobjects.NonEmptyString) string {
	return g.keyPrefix + g.AccountVerificationPrefix() + token.Value()
}

func (g *RedisKeyGenerator) PasswordResetKey(token valueobjects.NonEmptyString) string {
	return g.keyPrefix + g.PasswordResetPrefix() + token.Value()
}

func (g *RedisKeyGenerator) SessionKey(sessionId uuid.UUID) string {
	return g.keyPrefix + g.SessionPrefix() + sessionId.String()
}

func (g *RedisKeyGenerator) SessionsKey(userId userdomain.UserId) string {
	return g.keyPrefix + g.SessionsPrefix() + uuid.UUID(userId).String()
}

func NewRedisKeyGenerator(config *config.RedisConfig) *RedisKeyGenerator {
	return &RedisKeyGenerator{
		config.KeyPrefix,
	}
}
