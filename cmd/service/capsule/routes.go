package capsule

import (
	"net/http"

	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CapsuleStore
}

func NewHandler(store types.CapsuleStore) *Handler {
	return &Handler{store: store}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/capsules", handler.handleGetCapsules).Methods(http.MethodGet)
	router.HandleFunc("/capsules", handler.handleCreateCapsule).Methods(http.MethodPost)
}

func (handler *Handler) handleGetCapsules(w http.ResponseWriter, r *http.Request) {
	capsules, err := handler.store.GetCapsules(1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, capsules)
}

func (handler *Handler) handleCreateCapsule(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	// TODO: validate token
}
