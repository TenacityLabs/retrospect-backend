package user

import (
	"fmt"
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/cmd/service/auth"
	"github.com/TenacityLabs/time-capsule-backend/config"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
)

type Handler struct {
	userStore types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user", auth.WithJWTAuth(h.handleGetUser, h.userStore)).Methods("GET")
}

func (handler *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.LoginUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user exists
	user, err := handler.userStore.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid credentials"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (handler *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user already exists
	_, err = handler.userStore.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create user
	err = handler.userStore.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (handler *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())
	if userID == 0 {
		auth.PermissionDenied(w)
		return
	}

	user, err := handler.userStore.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]types.User{"user": *user})
}
