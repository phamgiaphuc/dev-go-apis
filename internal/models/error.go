package models

type APIError struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
	Stack   string `json:"stack,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) GetStack() string {
	return e.Stack
}

func (e *APIError) SetStack(stack string) {
	e.Stack = stack
}

func (e *APIError) WithStack(stack string) *APIError {
	e.Stack = stack
	return e
}
