package routes

type ResponseDTO[T any] struct {
	Data T `json:"data"`
}
