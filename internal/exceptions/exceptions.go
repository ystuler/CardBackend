package exceptions

const (
	ErrUserAlreadyExists   = "user already exists"
	ErrInvalidJSONFormat   = "invalid JSON format"
	ErrInternalServer      = "internal server error"
	ErrInvalidCredentials  = "invalid credentials"
	ErrInvalidToken        = "invalid token"
	ErrUnauthorized        = "unauthorized"
	ErrForbidden           = "forbidden"
	ErrMethodNotAllowed    = "method not allowed"
	ErrResourceNotFound    = "resource not found"
	ErrCollectionNotFound  = "collection not found"
	ErrCardNotFound        = "card not found"
	ErrValidationFailed    = "validation failed"
	ErrInvalidRequestBody  = "invalid request body"
	ErrInvalidCollectionID = "invalid collection ID, it must be an integer"
	ErrInvalidCardID       = "invalid card ID, it must be an integer"
)
