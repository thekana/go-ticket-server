package custom_error

type ValidationError struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type AuthorizationError struct {
	Code           uint64 `json:"code"`
	Message        string `json:"message"`
	HTTPStatusCode int    `json:"-"`
}

func (e *AuthorizationError) Error() string {
	return e.Message
}

type UserError struct {
	Code           uint64 `json:"code"`
	Message        string `json:"message"`
	HTTPStatusCode int    `json:"-"`
}

func (e *UserError) Error() string {
	return e.Message
}

type InternalError struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

func (e *InternalError) Error() string {
	return e.Message
}
