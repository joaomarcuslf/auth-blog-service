package types

type ErrorResponse struct {
	Error func() string
}
