package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
)

const (
	cursorSeparator string = "|"
	CursorMaxLength int    = 200
)

type KeysetCursor[TId fmt.Stringer] struct {
	Id        TId
	CreatedAt time.Time
}

func (c *KeysetCursor[TId]) String() string {
	fullValue := fmt.Sprintf(
		"%s%s%s",
		c.Id.String(),
		cursorSeparator,
		c.CreatedAt.Format(time.RFC3339Nano),
	)

	return base64.StdEncoding.EncodeToString([]byte(fullValue))
}

func NewKeysetCursor[TId fmt.Stringer](
	value string,
	parse func(string) (TId, error),
) (KeysetCursor[TId], error) {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) > CursorMaxLength {
		return KeysetCursor[TId]{}, domain.ValidationError(
			"cursor",
			[]string{"Invalid cursor"},
		)
	}

	decoded, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return KeysetCursor[TId]{}, domain.ValidationError(
			"cursor",
			[]string{"Invalid cursor"},
		)
	}

	parts := strings.Split(string(decoded), cursorSeparator)
	if len(parts) != 2 {
		return KeysetCursor[TId]{}, domain.ValidationError(
			"cursor",
			[]string{"Invalid cursor"},
		)
	}

	var validationErrors []string

	id, err := parse(parts[0])
	if err != nil {
		validationErrors = append(validationErrors, "Invalid cursor id format")
	}

	date, err := time.Parse(time.RFC3339Nano, parts[1])
	if err != nil {
		validationErrors = append(validationErrors, "Invalid cursor date format")
	}

	if len(validationErrors) > 0 {
		return KeysetCursor[TId]{}, domain.ValidationError(
			"cursor",
			validationErrors,
		)
	}

	return KeysetCursor[TId]{
		Id:        id,
		CreatedAt: date,
	}, nil
}
