package manager

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenManager interface {
	NewJWT(userID uuid.UUID, role string) (string, error)
	ParseJWT(tokenString string, field string) (string, error)
	ValidateJWT(tokenString string) (*jwt.MapClaims, error)
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewJWT(userID uuid.UUID, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["role"] = role
	claims["expired_time"] = time.Now().Add(15 * time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) ValidateJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, err
}

func (m *Manager) ParseJWT(tokenString string, field string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.signingKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims[field].(string), nil
	}

	return "", err
}
