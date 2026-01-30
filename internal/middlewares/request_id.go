package middlewares

import (
	"context"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/shared"
	"github.com/google/uuid"
)

func (m *Middlewares) RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, shared.RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(shared.RequestIDKey).(string); ok {
		return reqID
	}

	return ""
}
