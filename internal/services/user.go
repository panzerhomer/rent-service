package services

import (
	"avito/internal/domain"
	"avito/pkg/hasher"
	manager "avito/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
	GetByEmail(ctx context.Context, userEmail string) (domain.User, error)
	GetAll(ctx context.Context, offset int, limit int) ([]domain.User, error)
}

type UserService struct {
	log      *slog.Logger
	manager  manager.TokenManager
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository, manager manager.TokenManager, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		manager:  manager,
		log:      logger,
	}
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidUserType(userType string) bool {
	return userType == domain.Client || userType == domain.Moderator
}

func (s *UserService) Register(ctx context.Context, userRequest domain.UserRegisterRequest) (domain.UserRegisterResponse, error) {
	const op = "UserService.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", userRequest.Email),
	)

	log.Info("registering user")

	if !isValidUserType(userRequest.UserType) {
		s.log.Warn("invalid user role: " + userRequest.UserType)
		return domain.UserRegisterResponse{},
			fmt.Errorf("user usecase: register error: %w", domain.ErrUserBadType)
	}

	if !isValidEmail(userRequest.Email) {
		s.log.Warn("register error: invalid email", slog.String("email", userRequest.Email))
		return domain.UserRegisterResponse{},
			fmt.Errorf("register error: %w", domain.ErrUserBadEmail)
	}

	hashedPassword, err := hasher.HashPassword(userRequest.Password)
	if err != nil {
		s.log.Warn("hashing password failed: " + err.Error())
		return domain.UserRegisterResponse{},
			fmt.Errorf("hashing password failed: %v", err.Error())
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		s.log.Warn("generating uuid failed: " + err.Error())
		return domain.UserRegisterResponse{}, fmt.Errorf("generating uuid failed: %v", err.Error())
	}

	user := domain.User{
		ID:       uuid,
		Email:    userRequest.Email,
		Password: hashedPassword,
		Role:     userRequest.UserType,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		s.log.Warn("registration err: " + err.Error())
		if errors.Is(err, domain.ErrUserExist) {
			return domain.UserRegisterResponse{}, domain.ErrUserExist
		}
		return domain.UserRegisterResponse{},
			fmt.Errorf("registration err: %v", err.Error())
	}

	return domain.UserRegisterResponse{ID: uuid}, nil
}

func (s *UserService) Login(ctx context.Context, userRequest domain.UserLoginRequest) (domain.UserLoginResponse, error) {
	const op = "UserService.Login"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("loging user")

	savedUser, err := s.userRepo.GetByID(ctx, userRequest.ID)
	if err != nil {
		s.log.Warn("login error: " + err.Error())
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.UserLoginResponse{}, domain.ErrUserNotFound
		}
		return domain.UserLoginResponse{},
			fmt.Errorf("login error: %v", err.Error())
	}

	err = hasher.VerifyPassword(savedUser.Password, userRequest.Password)
	if err != nil {
		s.log.Warn("login error: " + err.Error())
		return domain.UserLoginResponse{},
			fmt.Errorf("login error: %v", err.Error())
	}

	token, err := s.manager.NewJWT(savedUser.ID, savedUser.Role)
	if err != nil {
		s.log.Warn("login error: " + err.Error())
		return domain.UserLoginResponse{},
			fmt.Errorf("user usecase: login error: %v", err.Error())
	}

	return domain.UserLoginResponse{Token: token}, nil
}

func (s *UserService) DummyLogin(ctx context.Context, userType string) (domain.UserLoginResponse, error) {
	const op = "UserService.DummyLogin"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("dummy loging user")

	if !isValidUserType(userType) {
		s.log.Warn("dummy login error: invalid role", slog.String("role", userType))
		return domain.UserLoginResponse{},
			fmt.Errorf("dummy login error: %w", domain.ErrUserBadType)
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		s.log.Warn("dummy login error: " + err.Error())
		return domain.UserLoginResponse{},
			fmt.Errorf("dummy login error: %v", err)
	}

	token, err := s.manager.NewJWT(uuid, userType)
	if err != nil {
		s.log.Warn("login error: " + err.Error())
		return domain.UserLoginResponse{},
			fmt.Errorf("login error: %v", err.Error())
	}

	return domain.UserLoginResponse{Token: token}, nil
}
