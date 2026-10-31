package constants

import (
	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

const (
	RedisAccountVerificationPrefix string = "account-verification:"
	RedisSessionPrefix             string = "session:"
	RedisSessionsPrefix            string = "sessions:"
	RedisPasswordResetPrefix       string = "password-reset:"
)

func RedisAccountVerificationKey(config *config.RedisConfig, token valueobjects.NonEmptyString) string {
	return config.KeyPrefix + RedisAccountVerificationPrefix + token.Value()
}

func RedisSessionKey(config *config.RedisConfig, sessionId uuid.UUID) string {
	return config.KeyPrefix + RedisSessionPrefix + sessionId.String()
}

func RedisSessionsKey(config *config.RedisConfig, userId userdomain.UserId) string {
	return config.KeyPrefix + RedisSessionsPrefix + uuid.UUID(userId).String()
}

func RedisPasswordResetKey(config *config.RedisConfig, userId userdomain.UserId) string {
	return config.KeyPrefix + RedisSessionsPrefix + uuid.UUID(userId).String()
}
