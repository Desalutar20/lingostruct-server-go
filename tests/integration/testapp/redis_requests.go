package testapp

import (
	"strings"
	"testing"
)

type TokenType int

const (
	AccountVerification TokenType = iota
	ResetPassword
	Session
)

func (a *TestApp) GetTokenFromRedis(t *testing.T, tokenType TokenType) *string {
	var prefix string

	switch tokenType {
	case AccountVerification:
		prefix = a.redisKeyGenerator.AccountVerificationPrefix()
	case ResetPassword:
		prefix = a.redisKeyGenerator.PasswordResetPrefix()
	case Session:
		prefix = a.redisKeyGenerator.SessionPrefix()
	}

	pattern := a.Config.Redis.KeyPrefix + prefix + "*"

	keys, err := a.redis.Keys(t.Context(), pattern).Result()
	if err != nil {
		t.Fatal()
	}

	if len(keys) == 0 {
		return nil
	}

	token := strings.Split(keys[0], prefix)[1]

	return &token
}
