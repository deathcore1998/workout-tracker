package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deathcore1998/workout-tracker/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
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
