package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"app/domain"
	"app/internal/controller/rest/user/converter"
	rest "app/internal/controller/rest/user/model"
	"app/internal/usecase/user"
	"app/pkg/logger"
)

type UserHandler struct {
	Service user.UserService
	logger  logger.Interface
}

func NewUserHandler(
	service user.UserService,
	logger logger.Interface,
	mux *http.ServeMux,
) (*UserHandler, error) {
	if service == nil {
		return nil, errors.New("UserHandler.NewUserHandler: service is null")
	}

	if logger == nil {
		return nil, errors.New("UserHandler.NewUserHandler: logger is null")
	}

	if mux == nil {
		return nil, errors.New("UserHandler.NewUserHandler: mux is null")
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
		h.logger.Error("UserHandler.CreateUser: error decoding user: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := h.Service.Create(r.Context(), converter.ToUserFromRest(user))
	if err != nil {
		h.logger.Error("UserHandler.CreateUser: error creating user: " + err.Error())
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
		if err == domain.ErrorUserNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			h.logger.Error("UserHandler.GetUser: error fetching user: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
		h.logger.Error("UserHandler.DeleteUser: error deleting user: " + err.Error())
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
