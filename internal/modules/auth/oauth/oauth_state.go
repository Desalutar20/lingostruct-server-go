package oauth

import (
	"fmt"
	"strings"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/google/uuid"
)

const separator string = "|"

type OAuthState struct {
	StateId         uuid.UUID
	AdditionalState *valueobjects.NonEmptyString
}

func (s OAuthState) String() string {
	if s.AdditionalState != nil {
		return fmt.Sprintf("%s%s%s", s.StateId.String(), separator, s.AdditionalState.Value())
	}

	return s.StateId.String()
}

func NewOAuthState(value string) (OAuthState, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return OAuthState{}, domain.FailureError("Invalid oauth state")
	}

	parts := strings.Split(value, separator)
	id, err := uuid.Parse(parts[0])
	if err != nil {
		return OAuthState{}, domain.FailureError("Invalid oauth state")
	}

	var additionalState *valueobjects.NonEmptyString

	if len(parts) == 2 {
		state, err := valueobjects.NewNonEmptyString(parts[1], "state", "State", 0, 0, nil)
		if err != nil {
			return OAuthState{}, domain.FailureError("Invalid oauth state")
		}

		additionalState = &state
	}

	return OAuthState{
		StateId:         id,
		AdditionalState: additionalState,
	}, nil
}
