package audio

import (
	"fmt"
	"net/http"

	"github.com/TenacityLabs/retrospect-backend/services/auth"
	"github.com/TenacityLabs/retrospect-backend/types"
	"github.com/TenacityLabs/retrospect-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	capsuleStore types.CapsuleStore
	userStore    types.UserStore
	fileStore    types.FileStore
	audioStore   types.AudioStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, fileStore types.FileStore, audioStore types.AudioStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		fileStore:    fileStore,
		audioStore:   audioStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/audios/create", auth.WithJWTAuth(handler.handleCreateAudio, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/audios/delete", auth.WithJWTAuth(handler.handleDeleteAudio, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateAudio(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateAudioPayload
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

	audioID, err := handler.audioStore.CreateAudio(userID, payload.CapsuleID, payload.ObjectName, payload.FileURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"id": audioID})
}

func (handler *Handler) handleDeleteAudio(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteAudioPayload
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

	objectName, err := handler.audioStore.DeleteAudio(userID, payload.CapsuleID, payload.AudioID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = handler.fileStore.DeleteFile(objectName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
