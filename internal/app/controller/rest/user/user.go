package user

import (
	"app/internal/app/controller/rest/user/converter"
	"app/internal/app/controller/rest/user/model"
	"app/internal/app/usecase/user"
	"app/internal/domain"
	"app/internal/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	Service user.Service
	logger  logger.Interface
}

func NewUserHandler(
	service user.Service,
	logger logger.Interface,
	mux *http.ServeMux,
) (*Handler, error) {
	if service == nil {
		return nil, errors.New("Handler.NewUserHandler: service is null")
	}

	if logger == nil {
		return nil, errors.New("Handler.NewUserHandler: logger is null")
	}

	if mux == nil {
		return nil, errors.New("Handler.NewUserHandler: mux is null")
	}

	handler := &Handler{Service: service, logger: logger}
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("GET /users", handler.GetUser)
	mux.HandleFunc("DELETE /users", handler.DeleteUser)
	return handler, nil
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u rest.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.logger.Error("Handler.CreateUser: error decoding user: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, err := h.Service.Create(r.Context(), converter.ToUserFromRest(u))
	if err != nil {
		h.logger.Error("Handler.CreateUser: error creating user: " + err.Error())
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
	if err != nil {
		h.logger.Error("Handler.CreateUser: error encoding user: " + err.Error())
		return
	}
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	u, err := h.Service.Get(r.Context(), uuid)
	if err != nil {
		if errors.Is(err, domain.ErrorUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			h.logger.Error("Handler.GetUser: error fetching user: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		h.logger.Error("Handler.CreateUser: error encoding user: " + err.Error())
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(r.Context(), uuid); err != nil {
		h.logger.Error("Handler.DeleteUser: error deleting user: " + err.Error())
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
