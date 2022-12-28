package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

type JWTCreator struct {
	secretKey string
}

func NewJWTCreator(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("Invalid secret key length: must be %d characters", minSecretKeySize)
	}

	return &JWTCreator{secretKey}, nil
}

func (maker *JWTCreator) CreateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwt.SignedString([]byte(maker.secretKey))
}

func (maker *JWTCreator) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenError
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ExpiredTokenErr) {
			return nil, ExpiredTokenErr
		}
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenError
	}

	return payload, nil
}
