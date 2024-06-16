package capsule

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TenacityLabs/time-capsule-backend/services/auth"
	"github.com/TenacityLabs/time-capsule-backend/types"
	"github.com/TenacityLabs/time-capsule-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	capsuleStore        types.CapsuleStore
	userStore           types.UserStore
	songStore           types.SongStore
	questionAnswerStore types.QuestionAnswerStore
	writingStore        types.WritingStore
	photoStore          types.PhotoStore
	audioStore          types.AudioStore
}

func NewHandler(
	capsuleStore types.CapsuleStore,
	userStore types.UserStore,
	songStore types.SongStore,
	questionAnswerStore types.QuestionAnswerStore,
	writingStore types.WritingStore,
	photoStore types.PhotoStore,
	audioStore types.AudioStore,
) *Handler {
	return &Handler{
		capsuleStore:        capsuleStore,
		userStore:           userStore,
		songStore:           songStore,
		questionAnswerStore: questionAnswerStore,
		writingStore:        writingStore,
		photoStore:          photoStore,
		audioStore:          audioStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/capsules", auth.WithJWTAuth(handler.handleGetCapsules, handler.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/capsules/get-by-id/{capsuleId}", auth.WithJWTAuth(handler.handleGetCapsuleById, handler.userStore)).Methods(http.MethodGet)
	router.HandleFunc("/capsules/create", auth.WithJWTAuth(handler.handleCreateCapsule, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/capsules/join", auth.WithJWTAuth(handler.handleJoinCapsule, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/capsules/delete", auth.WithJWTAuth(handler.handleDeleteCapsule, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleGetCapsules(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())

	capsules, err := handler.capsuleStore.GetCapsules(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, capsules)
}

func (handler *Handler) handleGetCapsuleById(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIdFromContext(r.Context())
	vars := mux.Vars(r)
	capsuleIdStr, ok := vars["capsuleId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("capsuleId not provided"))
		return
	}
	capsuleId, err := strconv.Atoi(capsuleIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid capsuleId"))
		return
	}

	capsule, err := handler.capsuleStore.GetCapsuleById(userID, uint(capsuleId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	songs, err := handler.songStore.GetSongs(uint(capsule.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	questionAnswers, err := handler.questionAnswerStore.GetQuestionAnswers(uint(capsule.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	writings, err := handler.writingStore.GetWritings(uint(capsule.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	photos, err := handler.photoStore.GetPhotos(uint(capsule.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	audios, err := handler.audioStore.GetAudios(uint(capsule.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: get all capsule components

	utils.WriteJSON(w, http.StatusOK, types.GetCapsuleByIdResponse{
		Capsule:         capsule,
		Songs:           songs,
		QuestionAnswers: questionAnswers,
		Writings:        writings,
		Photos:          photos,
		Audios:          audios,
	})
}

func (handler *Handler) handleCreateCapsule(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateCapsulePayload
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

	capsuleID, err := handler.capsuleStore.CreateCapsule(userID, payload.Vessel, payload.Public)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"capsuleId": capsuleID})
}

func (handler *Handler) handleJoinCapsule(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.JoinCapsulePayload
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

	err = handler.capsuleStore.JoinCapsule(userID, payload.Code)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (handler *Handler) handleDeleteCapsule(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteCapsulePayload
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

	err = handler.capsuleStore.DeleteCapsule(userID, payload.CapsuleID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}
