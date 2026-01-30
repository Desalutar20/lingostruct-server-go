package shared

type ContextKey string

const (
	UserContextKey ContextKey = "user"
	RequestIDKey   ContextKey = "request_id"
)
