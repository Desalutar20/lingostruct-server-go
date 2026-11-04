package oauth

import (
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/google/uuid"
)

func TestOAuthState(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		stateID := uuid.New()

		t.Run("without additional state", func(t *testing.T) {
			state := OAuthState{
				StateId: stateID,
			}

			got := state.String()
			want := stateID.String()

			if got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}
		})

		t.Run("with additional state", func(t *testing.T) {
			additionalState := "dashboard"

			state := OAuthState{
				StateId:         stateID,
				AdditionalState: stringToNonEmptyString(t, additionalState),
			}

			got := state.String()
			want := stateID.String() + separator + additionalState

			if got != want {
				t.Errorf("String() = %q, want %q", got, want)
			}
		})
	})

	t.Run("NewOAuthState", func(t *testing.T) {
		stateID := uuid.New()

		tests := []struct {
			name              string
			value             string
			wantStateID       uuid.UUID
			wantAdditional    string
			wantHasAdditional bool
			wantErr           bool
		}{
			{
				name:        "valid state without additional state",
				value:       stateID.String(),
				wantStateID: stateID,
			},
			{
				name:              "valid state with additional state",
				value:             stateID.String() + separator + "dashboard",
				wantStateID:       stateID,
				wantAdditional:    "dashboard",
				wantHasAdditional: true,
			},
			{
				name:              "additional state with surrounding whitespace",
				value:             stateID.String() + separator + "  dashboard  ",
				wantStateID:       stateID,
				wantAdditional:    "dashboard",
				wantHasAdditional: true,
			},
			{
				name:    "empty value",
				value:   "",
				wantErr: true,
			},
			{
				name:    "whitespace only",
				value:   "   ",
				wantErr: true,
			},
			{
				name:    "invalid UUID",
				value:   "invalid-uuid",
				wantErr: true,
			},
			{
				name:    "empty additional state",
				value:   stateID.String() + separator,
				wantErr: true,
			},
			{
				name:    "whitespace only additional state",
				value:   stateID.String() + separator + "   ",
				wantErr: true,
			},
			{
				name:              "multiple separators",
				value:             stateID.String() + separator + "dashboard" + separator + "settings",
				wantStateID:       stateID,
				wantHasAdditional: false,
			},
			{
				name:        "surrounding whitespace",
				value:       "  " + stateID.String() + "  ",
				wantStateID: stateID,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := NewOAuthState(tt.value)

				if tt.wantErr {
					if err == nil {
						t.Fatal("NewOAuthState() expected an error, got nil")
					}
					return
				}

				if err != nil {
					t.Fatalf("NewOAuthState() unexpected error: %v", err)
				}

				if got.StateId != tt.wantStateID {
					t.Errorf("StateId = %v, want %v", got.StateId, tt.wantStateID)
				}

				if tt.wantHasAdditional {
					if got.AdditionalState == nil {
						t.Fatal("AdditionalState = nil, want non-nil")
					}

					if got.AdditionalState.Value() != tt.wantAdditional {
						t.Errorf(
							"AdditionalState = %q, want %q",
							got.AdditionalState.Value(),
							tt.wantAdditional,
						)
					}
				} else if got.AdditionalState != nil {
					t.Errorf(
						"AdditionalState = %q, want nil",
						got.AdditionalState.Value(),
					)
				}
			})
		}
	})
}

func stringToNonEmptyString(t *testing.T, value string) *valueobjects.NonEmptyString {
	t.Helper()

	result, err := valueobjects.NewNonEmptyString(
		value,
		"state",
		"State",
		0,
		0,
		nil,
	)
	if err != nil {
		t.Fatalf("NewNonEmptyString() unexpected error: %v", err)
	}

	return &result
}
