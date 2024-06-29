package miscFile

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
	capsuleStore  types.CapsuleStore
	userStore     types.UserStore
	fileStore     types.FileStore
	miscFileStore types.MiscFileStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, fileStore types.FileStore, miscFileStore types.MiscFileStore) *Handler {
	return &Handler{
		capsuleStore:  capsuleStore,
		userStore:     userStore,
		fileStore:     fileStore,
		miscFileStore: miscFileStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/misc-files/create", auth.WithJWTAuth(handler.handleCreateMiscFile, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/misc-files/delete", auth.WithJWTAuth(handler.handleDeleteMiscFile, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateMiscFile(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateMiscFilePayload
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

	miscFileID, err := handler.miscFileStore.CreateMiscFile(userID, payload.CapsuleID, payload.ObjectName, payload.FileURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"miscFileId": miscFileID})
}

func (handler *Handler) handleDeleteMiscFile(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteMiscFilePayload
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

	objectName, err := handler.miscFileStore.DeleteMiscFile(userID, payload.CapsuleID, payload.MiscFileID)
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
