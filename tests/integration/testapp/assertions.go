package testapp

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
)

func (a *TestApp) AssertValidStatus(t *testing.T, resp *http.Response, statusCode int) {
	t.Helper()

	if resp.StatusCode != statusCode {
		t.Fatalf(
			"expected status %d, got %d",
			statusCode,
			resp.StatusCode,
		)
	}
}

func (a *TestApp) AssertValidationError(t *testing.T, resp *http.Response, field string) {
	var errorDTO routes.ErrorDTO
	if err := json.NewDecoder(resp.Body).Decode(&errorDTO); err != nil {
		t.Fatalf("decoding error response: %v", err)
	}

	if errorDTO.Errors == nil {
		t.Fatal("expected validation errors")
	}

	if _, ok := errorDTO.Errors[field]; !ok {
		t.Fatalf("expected validation error for field %q, got %v", field, errorDTO.Errors)
	}
}
