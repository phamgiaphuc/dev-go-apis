package models

type APIError struct {
	Code    int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Stack   string      `json:"stack,omitempty"`
	Errors  interface{} `json:"errors"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) SetMessage(message string) {
	e.Message = message
}

func (e *APIError) SetStack(stack string) {
	e.Stack = stack
}

func (e *APIError) SetErrors(errors interface{}) {
	e.Errors = errors
}

func (e *APIError) GetError() *APIError {
	return e
}

func (e *APIError) WithStack(stack string) *APIError {
	e.Stack = stack
	return e
}
