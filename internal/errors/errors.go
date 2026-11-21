package errors

type JSONError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
