// API response using generics
package response

type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    *T     `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func Success[T any](data T) APIResponse[T] {
	return APIResponse[T]{Success: true, Data: &data}
}

func Failure[T any](errMsg string) APIResponse[T] {
	return APIResponse[T]{Success: false, Error: errMsg}
}
