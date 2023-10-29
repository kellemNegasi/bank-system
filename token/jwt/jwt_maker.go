package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeyLength = 32

type JWTMaker struct {
	secretKey string
}

func NewJWT(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("Invalid secret key length. Secret key should be %v long", minSecretKeyLength)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a token for the username and specified duration.
func (jwtMaker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payLoad, err := NewPayLoad(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)
	return jwtToken.SignedString([]byte(jwtMaker.secretKey))

}

// Verify returns the payload contained in the given token.
func (jwtMaker *JWTMaker) Verify(token string) (*JWTPayload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(jwtMaker.secretKey), nil

	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyfunc)
	if err != nil {
		// TODO: differentiate different types of errors.
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*JWTPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
