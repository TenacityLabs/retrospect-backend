package writing

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
	writingStore types.WritingStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, writingStore types.WritingStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		writingStore: writingStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/writings/create", auth.WithJWTAuth(handler.handleCreateWriting, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/writings/delete", auth.WithJWTAuth(handler.handleDeleteWriting, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateWriting(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateWritingPayload
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

	// check if user is member of capsule
	_, err = handler.capsuleStore.GetCapsuleById(userID, payload.CapsuleID)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("could not find capsule with id %d", payload.CapsuleID))
		return
	}

	writingID, err := handler.writingStore.CreateWriting(userID, payload.CapsuleID, payload.Writing)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"writingId": writingID})
}

func (handler *Handler) handleDeleteWriting(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteWritingPayload
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

	err = handler.writingStore.DeleteWriting(userID, payload.CapsuleID, payload.WritingID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}
