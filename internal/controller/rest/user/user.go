package user

import (
	"app/domain"
	"app/internal/controller/rest/user/converter"
	rest "app/internal/controller/rest/user/model"
	"app/internal/usecase/user"
	"app/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"
)

type UserHandler struct {
	Service user.UserService
	logger  logger.Interface
}

func NewUserHandler(service user.UserService, logger logger.Interface, mux *http.ServeMux) (*UserHandler, error) {

	if service == nil {
		return nil, errors.New("service is null")
	}

	if logger == nil {
		return nil, errors.New("logger is null")
	}

	if mux == nil {
		return nil, errors.New("mux is null")
	}

	handler := &UserHandler{Service: service, logger: logger}
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("GET /users", handler.GetUser)
	mux.HandleFunc("DELETE /users", handler.DeleteUser)
	return handler, nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user rest.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Error decoding user: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := h.Service.Create(r.Context(), converter.ToUserFromRest(user))
	if err != nil {
		h.logger.Error("Error creating user: " + err.Error())
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	user, err := h.Service.Get(r.Context(), uuid)
	if err != nil {
		h.logger.Error("Error fetching user: " + err.Error())
		if err == domain.ErrorUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Could not fetch user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(r.Context(), uuid); err != nil {
		h.logger.Error("Error deleting user: " + err.Error())
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
