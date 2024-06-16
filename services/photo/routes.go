package photo

import (
	"fmt"
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/services/auth"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	capsuleStore types.CapsuleStore
	userStore    types.UserStore
	fileStore    types.FileStore
	photoStore   types.PhotoStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, fileStore types.FileStore, photoStore types.PhotoStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		fileStore:    fileStore,
		photoStore:   photoStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/photos/create", auth.WithJWTAuth(handler.handleCreatePhoto, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/photos/delete", auth.WithJWTAuth(handler.handleDeletePhoto, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreatePhoto(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreatePhotoPayload
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

	photoID, err := handler.photoStore.CreatePhoto(userID, payload.CapsuleID, payload.ObjectName, payload.FileURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"photoId": photoID})
}

func (handler *Handler) handleDeletePhoto(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeletePhotoPayload
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

	objectName, err := handler.photoStore.DeletePhoto(userID, payload.CapsuleID, payload.PhotoID)
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
