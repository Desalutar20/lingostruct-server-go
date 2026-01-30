package setup

import (
	"strings"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/constants"
)

func (a *TestApp) GetVerificationToken(t *testing.T) string {
	keys, err := a.redis.Keys(t.Context(), "*").Result()
	if err != nil {
		t.Fatalf("getVerificationToken failed: %v", err)
	}

	verificationToken := ""

	for i := len(keys) - 1; i >= 0; i-- {
		if strings.HasPrefix(keys[i], constants.RedisVerificationPrefix) {
			verificationToken = strings.Split(keys[i], constants.RedisVerificationPrefix)[1]
			break
		}
	}

	return verificationToken
}

func (a *TestApp) GetResetPasswordToken(t *testing.T) string {
	keys, err := a.redis.Keys(t.Context(), "*").Result()
	if err != nil {
		t.Fatalf("getResetPasswordToken failed: %v", err)
	}

	resetPasswordToken := ""

	for i := len(keys) - 1; i >= 0; i-- {
		if strings.HasPrefix(keys[i], constants.RedisResetPasswordPrefix) {
			resetPasswordToken = strings.Split(keys[i], constants.RedisResetPasswordPrefix)[1]
			break
		}
	}

	return resetPasswordToken
}

func (a *TestApp) ClearCache(t *testing.T) {
	_, err := a.redis.FlushAll(t.Context()).Result()
	if err != nil {
		t.Fatalf("clearCache failed: %v", err)
	}
}
