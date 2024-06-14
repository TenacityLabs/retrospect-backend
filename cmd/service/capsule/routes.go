package capsule

import (
	"fmt"
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/cmd/service/auth"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	capsuleStore types.CapsuleStore
	userStore    types.UserStore
	songStore    types.SongStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, songStore types.SongStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		songStore:    songStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/capsules", auth.WithJWTAuth(handler.handleGetCapsules, handler.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/capsules/create", auth.WithJWTAuth(handler.handleCreateCapsule, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/capsules/join", auth.WithJWTAuth(handler.handleJoinCapsule, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleGetCapsules(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())

	capsules, err := handler.capsuleStore.GetCapsules(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, capsules)
}

func (handler *Handler) handleCreateCapsule(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateCapsulePayload
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

	userID := auth.GetUserIdFromContext(r.Context())

	capsuleID, err := handler.capsuleStore.CreateCapsule(userID, payload.Vessel, payload.Public)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"capsuleId": capsuleID})
}

func (handler *Handler) handleJoinCapsule(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.JoinCapsulePayload
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

	userID := auth.GetUserIdFromContext(r.Context())

	err = handler.capsuleStore.JoinCapsule(userID, payload.Code)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}
