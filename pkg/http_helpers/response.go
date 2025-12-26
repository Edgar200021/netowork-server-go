package http_helpers

type SuccessResponse[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type ValidationErrorResponse struct {
	Status int                 `json:"status"`
	Errors map[string][]string `json:"errors"`
}
