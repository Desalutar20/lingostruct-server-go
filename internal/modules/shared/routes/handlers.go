package routes

import (
	"net/http"
)

func RateLimitHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, http.StatusTooManyRequests, ErrorDTO{
		Error: "Too many requests",
		Code:  "TOO_MANY_REQUESTS",
	})
}
