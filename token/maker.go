package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	//CreateToken creates a new Token for a specific username and valid duration
	CreateToken(userName string, duration time.Duration) (string, error)
	//VerifyToken check if token is valid or not
	VerifyJWTToken(token string) (*Payload, error)
}
