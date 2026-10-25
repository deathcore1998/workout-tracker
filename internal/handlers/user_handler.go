package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deathcore1998/workout-tracker/internal/auth"
	"github.com/deathcore1998/workout-tracker/internal/services"
)

type UserHandler struct {
	service   *services.UserService
	jwtSecret string
}

func NewUserHandler(service *services.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{service: service, jwtSecret: jwtSecret}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(writer http.ResponseWriter, clientData *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(clientData.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(clientData.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(user)
}

func (h *UserHandler) Login(writer http.ResponseWriter, clientData *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(clientData.Body).Decode(&req); err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(clientData.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(writer, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, h.jwtSecret)
	if err != nil {
		http.Error(writer, "failed to generate token", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(LoginResponse{Token: token})
}
