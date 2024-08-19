package domain

import (
	"errors"

	"github.com/google/uuid"
)

const (
	Client    = "client"
	Moderator = "moderator"
)

var (
	ErrUserBadType     = errors.New("bad user type")
	ErrUserBadRequest  = errors.New("bad empty request")
	ErrUserBadEmail    = errors.New("bad mail")
	ErrUserBadPassword = errors.New("bad password")
	ErrUserBadID       = errors.New("bad id")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExist    = errors.New("user already exist")
)

type User struct {
	ID       uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

func (u *UserRegisterRequest) Validate() error {
	if u.Email == "" {
		return ErrUserBadEmail
	}
	if u.UserType == "" {
		return ErrUserBadType
	}
	if u.Password == "" || len(u.Password) < 4 {
		return ErrUserBadPassword
	}
	return nil
}

type UserRegisterResponse struct {
	ID uuid.UUID `json:"user_id"`
}

type UserLoginRequest struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

func (u *UserLoginRequest) Validate() error {
	if u.ID.String() == "" {
		return ErrUserBadID
	}
	if u.Password == "" || len(u.Password) < 4 {
		return ErrUserBadPassword
	}
	return nil
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type DummyLoginRequest struct {
	UserType string `json:"user_type"`
}
