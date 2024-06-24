package file

import (
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/config"
	"github.com/TenacityLabs/time-capsule-backend/services/auth"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
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
