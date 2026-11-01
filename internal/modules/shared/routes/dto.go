package routes

type ResponseDTO[T any] struct {
	Data T `json:"data"`
}

type ErrorDTO struct {
	Error  string              `json:"error"`
	Errors map[string][]string `json:"errors,omitempty"`
}
