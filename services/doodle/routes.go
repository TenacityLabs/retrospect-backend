package doodle

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
	doodleStore  types.DoodleStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, fileStore types.FileStore, doodleStore types.DoodleStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		fileStore:    fileStore,
		doodleStore:  doodleStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/doodles/create", auth.WithJWTAuth(handler.handleCreateDoodle, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/doodles/delete", auth.WithJWTAuth(handler.handleDeleteDoodle, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateDoodle(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateDoodlePayload
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

	doodleID, err := handler.doodleStore.CreateDoodle(userID, payload.CapsuleID, payload.ObjectName, payload.FileURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"id": doodleID})
}

func (handler *Handler) handleDeleteDoodle(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteDoodlePayload
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

	objectName, err := handler.doodleStore.DeleteDoodle(userID, payload.CapsuleID, payload.DoodleID)
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
