package handlers

import (
	"avito/internal/domain"
	manager "avito/pkg/jwt"
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"io"
	"net/http"
	"time"
)

type UserService interface {
	Register(ctx context.Context, userRequest domain.UserRegisterRequest) (domain.UserRegisterResponse, error)
	Login(ctx context.Context, userRequest domain.UserLoginRequest) (domain.UserLoginResponse, error)
	DummyLogin(ctx context.Context, userType string) (domain.UserLoginResponse, error)
}

type UserHandler struct {
	userService UserService
	manager     manager.TokenManager
	log         *slog.Logger
}

func NewUserHandler(userService UserService, manager manager.TokenManager, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         logger,
		manager:     manager,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Register"

	log := h.log.With(
		slog.String("op", op),
	)

	log.Info("registering user")

	var (
		respBody        []byte
		registerRequest domain.UserRegisterRequest
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Warn("register error: " + err.Error())
		respBody := CreateErrorResponse(r.Context(), ErrorReadHTTPBody, "can't read body")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		h.log.Warn("register error: " + err.Error())
		respBody = CreateErrorResponse(r.Context(), ErrorUnmarshalHTTPBody, "can't unmarshal request")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	if registerRequest.Email == "" || registerRequest.Password == "" || registerRequest.UserType == "" {
		h.log.Warn("register error: some filed empty")
		respBody = CreateErrorResponse(r.Context(), ErrorRegisterUser, "can't register user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBody)
		return
	}

	ctx, cancelRegistartion := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelRegistartion()

	userRegisterResponse, err := h.userService.Register(ctx, registerRequest)
	if err != nil {
		h.log.Warn("register error: " + err.Error())
		respBody = CreateErrorResponse(r.Context(), ErrorRegisterUser, ErrorRegisterUserMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	respBody, err = json.Marshal(userRegisterResponse)
	if err != nil {
		h.log.Warn("register error: " + err.Error())
		respBody = CreateErrorResponse(r.Context(), ErrorMarshalHTTPBody, ErrorMarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Login"

	log := h.log.With(
		slog.String("op", op),
	)

	log.Info("loging user")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Warn("login error: " + err.Error())
		respBody := CreateErrorResponse(r.Context(), ErrorReadHTTPBody, ReadHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	var loginRequest domain.UserLoginRequest
	if err = json.Unmarshal(body, &loginRequest); err != nil {
		h.log.Warn("login error: " + err.Error())
		respBody := CreateErrorResponse(r.Context(), ErrorUnmarshalHTTPBody, UnmarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	if err := loginRequest.Validate(); err != nil {
		h.log.Warn("login error: empty ID or password")
		respBody := CreateErrorResponse(r.Context(), ErrorLoginUser, ErrorLoginUserMsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(respBody)
		return
	}

	ctx, cancelLogin := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelLogin()

	loginResponse, err := h.userService.Login(ctx, loginRequest)
	if err != nil {
		h.log.Warn("login error: " + err.Error())

		var statusCode int
		var messageError string
		if errors.Is(err, domain.ErrUserNotFound) {
			statusCode = http.StatusNotFound
			messageError = domain.ErrUserNotFound.Error()
		} else {
			statusCode = http.StatusInternalServerError
			messageError = ErrorLoginUserMsg
		}

		w.WriteHeader(statusCode)
		respBody := CreateErrorResponse(ctx, ErrorLoginUser, messageError)
		w.Write(respBody)
		return
	}

	respBody, err := json.Marshal(loginResponse)
	if err != nil {
		h.log.Warn("login error: " + err.Error())
		respBody = CreateErrorResponse(r.Context(), ErrorMarshalHTTPBody, ErrorMarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    loginResponse.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (h *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	userType := r.URL.Query().Get("user_type")
	if userType == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancelLogin := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelLogin()

	dummyLoginResponse, err := h.userService.DummyLogin(ctx, userType)
	if err != nil {
		h.log.Warn("dummy login error: " + err.Error())
		respBody := CreateErrorResponse(r.Context(), ErrorDummyLogin, "can't dummy login")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	respBody, err := json.Marshal(dummyLoginResponse)
	if err != nil {
		h.log.Warn("dummy login error: " + err.Error())
		respBody = CreateErrorResponse(r.Context(), ErrorMarshalHTTPBody, "can't marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    dummyLoginResponse.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
