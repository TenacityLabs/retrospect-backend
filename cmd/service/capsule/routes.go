package capsule

import (
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/cmd/service/auth"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	capsuleStore types.CapsuleStore
	userStore    types.UserStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/capsules", auth.WithJWTAuth(handler.handleGetCapsules, handler.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/capsules/create", auth.WithJWTAuth(handler.handleCreateCapsule, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleGetCapsules(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())

	capsules, err := handler.capsuleStore.GetCapsules(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, capsules)
}

func (handler *Handler) handleCreateCapsule(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())

	capsuleID, err := handler.capsuleStore.CreateCapsule(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"capsuleId": capsuleID})
}
