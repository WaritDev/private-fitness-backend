package repositories

type JWTPayload struct {
	Sub       string `json:"sub"`
	Role      string `json:"role"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type AuthRepo interface {
	// SignJWT creates a JWT token with 7 days expiration
	SignJWT(payload JWTPayload) (string, error)

	// VerifyJWT verifies and decodes a JWT token
	VerifyJWT(token string) (*JWTPayload, error)
}
