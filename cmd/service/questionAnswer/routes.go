package questionAnswer

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
	songStore    types.SongStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, songStore types.SongStore) *Handler {
	return &Handler{
		capsuleStore: capsuleStore,
		userStore:    userStore,
		songStore:    songStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/songs/create", auth.WithJWTAuth(handler.handleCreateSong, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateSong(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateSongPayload
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

	songID, err := handler.songStore.CreateSong(userID, payload.CapsuleID, payload.SpotifyID, payload.Name, payload.ArtistName, payload.AlbumArtURL)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"songID": songID})
}
