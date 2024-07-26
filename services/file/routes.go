package file

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/TenacityLabs/retrospect-backend/config"
	"github.com/TenacityLabs/retrospect-backend/services/auth"
	"github.com/TenacityLabs/retrospect-backend/types"
	"github.com/TenacityLabs/retrospect-backend/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	userStore types.UserStore
	fileStore types.FileStore
}

func NewHandler(userStore types.UserStore, fileStore types.FileStore) *Handler {
	return &Handler{
		userStore: userStore,
		fileStore: fileStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/files/upload", auth.WithJWTAuth(handler.handleFileUpload, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/files/update", auth.WithJWTAuth(handler.handleFileUpdate, handler.userStore)).Methods(http.MethodPost)

	if config.Envs.GoEnv == "development" {
		router.HandleFunc("/files/delete", auth.WithJWTAuth(handler.handleFileDelete, handler.userStore)).Methods(http.MethodPost)
	}
}

func (handler *Handler) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Set a max memory limit of 10MB for parsing
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := auth.GetUserIdFromContext(r.Context())

	objectName, fileURL, err := handler.fileStore.UploadFile(userID, file, fileHeader)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"objectName": objectName, "fileURL": fileURL})
}

func (handler *Handler) handleFileUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Set a max memory limit of 10MB for parsing
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	objectName := r.FormValue("objectName")
	if objectName == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("objectName is required"))
		return
	}
	userID := auth.GetUserIdFromContext(r.Context())

	if !strings.HasPrefix(objectName, fmt.Sprintf("%d-", userID)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you are not allowed to update this file, only the creator of the file can update it"))
		return
	}

	err = handler.fileStore.DeleteFile(objectName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, _, err = handler.fileStore.UploadFileWithName(objectName, file, fileHeader)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (handler *Handler) handleFileDelete(w http.ResponseWriter, r *http.Request) {
	var payload types.DeleteFilePayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = handler.fileStore.DeleteFile(payload.ObjectName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
