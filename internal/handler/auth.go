package handler

import (
	"encoding/json"
	"errors"
	"helpdesk-api/internal/model"
	"helpdesk-api/internal/service"
	"net/http"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Login godoc
// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}
	if user.Email == "" || user.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required", nil)
		return
	}

	jwtToken, err := h.userService.Login(r.Context(), user.Email, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid email or password", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "login failed", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": jwtToken})

}

// Register godoc
// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if user.Email == "" || user.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required", nil)
		return
	}

	err = h.userService.Register(r.Context(), &user)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			writeError(w, http.StatusConflict, "email already exists", nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "registration failed", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User registered successfully"})

}
