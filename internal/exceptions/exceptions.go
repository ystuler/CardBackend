package exceptions

const (
	ErrUserAlreadyExists   = "user already exists"
	ErrInvalidJSONFormat   = "invalid JSON format"
	ErrInternalServer      = "internal server error"
	ErrInvalidCredentials  = "invalid credentials"
	ErrInvalidToken        = "invalid token"
	ErrInvalidCollectionID = "invalid collection ID, it must be an integer"
	ErrInvalidCardID       = "invalid card ID, it must be an integer"
)
