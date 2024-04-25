package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeySize = 32
)

var (
	ErrInvalidKeySize = fmt.Errorf("inalid key size: must have atleadt %d characters", minSecretKeySize)
	ErrInvalidToken   = errors.New("token is invalid")
	ErrTokenExpired   = errors.New("token has expired")
)

type Claims struct {
	JWTPayload
	jwt.RegisteredClaims
}

// JWTMaker is a new JSON web token maker
type JwtMaker struct {
	secretKey []byte
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, ErrInvalidKeySize
	}
	return &JwtMaker{[]byte(secretkey)}, nil
}

func (maker *JwtMaker) CreateToken(userName string, duration time.Duration) (string, error) {
	jwtpayload, err := NewJWTPayload(userName, duration)
	if err != nil {
		return "", err
	}
	claims := &Claims{
		JWTPayload: *jwtpayload,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(maker.secretKey)

}

//VerifyToken check if token is valid or not

func (maker *JwtMaker) VerifyJWTToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			{
				return nil, ErrTokenExpired
			}
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			{
				return nil, ErrInvalidToken
			}
		}
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}
	jwtPayload := jwtToken.Claims.(*JWTPayload)

	return &jwtPayload.Payload, err
}
