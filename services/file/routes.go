package file

import (
	"log"
	"net/http"

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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	// router.HandleFunc("/files/upload", auth.WithJWTAuth(handler.handleFileUpload, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/files/upload", handler.handleFileUpload).Methods(http.MethodPost)
}

func (handler *Handler) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// get json payload
	// var payload types.UploadFilePayload
	enableCors(&w)

	log.Printf("request: %v", r)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"fileURL": ""})
}
