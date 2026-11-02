package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/deathcore1998/workout-tracker/internal/middleware"
	"github.com/deathcore1998/workout-tracker/internal/models"
	"github.com/deathcore1998/workout-tracker/internal/services"
)

type WorkoutHandler struct {
	service *services.WorkoutService
}

func NewWorkoutHandler(service *services.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{service: service}
}

type CreateWorkoutRequest struct {
	Name  string    `json:"name"`
	Date  time.Time `json:"date"`
	Notes string    `json:"notes"`
}

type UpdateWorkoutRequest = CreateWorkoutRequest

func (handler *WorkoutHandler) Create(writer http.ResponseWriter, clientData *http.Request) {
	userID, ok := clientData.Context().Value(middleware.UserID).(int)
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request CreateWorkoutRequest
	err := json.NewDecoder(clientData.Body).Decode(&request)
	if err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	workout := &models.Workout{
		UserID: userID,
		Name:   request.Name,
		Date:   request.Date,
		Notes:  request.Notes,
	}

	err = handler.service.Create(clientData.Context(), workout)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(workout)
}

func (handler *WorkoutHandler) List(writer http.ResponseWriter, clientData *http.Request) {
	userID, ok := clientData.Context().Value(middleware.UserID).(int)
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	workouts, err := handler.service.GetByUserID(clientData.Context(), userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(workouts)
}

func (handler *WorkoutHandler) Get(writer http.ResponseWriter, clientData *http.Request) {
	userID, ok := clientData.Context().Value(middleware.UserID).(int)
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := clientData.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Invalid id", http.StatusBadRequest)
		return
	}

	workout, err := handler.service.GetByID(clientData.Context(), id, userID)
	if err != nil {
		http.Error(writer, "workout not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(workout)
}

func (handler *WorkoutHandler) Update(writer http.ResponseWriter, clientData *http.Request) {
	userID, ok := clientData.Context().Value(middleware.UserID).(int)
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := clientData.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Invalid id", http.StatusBadRequest)
		return
	}

	var request UpdateWorkoutRequest
	err = json.NewDecoder(clientData.Body).Decode(&request)
	if err != nil {
		http.Error(writer, "invalid request body", http.StatusBadRequest)
		return
	}

	workout := &models.Workout{
		ID:     id,
		UserID: userID,
		Name:   request.Name,
		Date:   request.Date,
		Notes:  request.Notes,
	}

	err = handler.service.Update(clientData.Context(), workout)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(workout)
}

func (handler *WorkoutHandler) Delete(writer http.ResponseWriter, clientData *http.Request) {
	userID, ok := clientData.Context().Value(middleware.UserID).(int)
	if !ok {
		http.Error(writer, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := clientData.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(writer, "Invalid id", http.StatusBadRequest)
		return
	}

	err = handler.service.Delete(clientData.Context(), id, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
