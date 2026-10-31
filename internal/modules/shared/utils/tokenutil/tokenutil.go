package tokenutil

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
)

func Generate(length int) (valueobjects.NonEmptyString, error) {
	if length <= 0 {
		return valueobjects.NonEmptyString{}, domain.InternalError("token length must be positive")
	}

	bytes := make([]byte, (length+1)/2)

	if _, err := rand.Read(bytes); err != nil {
		return valueobjects.NonEmptyString{}, err
	}

	hex := hex.EncodeToString(bytes)[:length]

	vo, err := valueobjects.NewNonEmptyString(hex, "token", "Token", length, 0, nil)
	if err != nil {
		return valueobjects.NonEmptyString{}, err
	}

	return vo, nil
}
