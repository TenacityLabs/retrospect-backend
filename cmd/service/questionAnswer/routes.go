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
	capsuleStore        types.CapsuleStore
	userStore           types.UserStore
	questionAnswerStore types.QuestionAnswerStore
}

func NewHandler(capsuleStore types.CapsuleStore, userStore types.UserStore, questionAnswerStore types.QuestionAnswerStore) *Handler {
	return &Handler{
		capsuleStore:        capsuleStore,
		userStore:           userStore,
		questionAnswerStore: questionAnswerStore,
	}
}

func (handler *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/question-answers/create", auth.WithJWTAuth(handler.handleCreateQuestionAnswer, handler.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/question-answers/delete", auth.WithJWTAuth(handler.handleDeleteQuestionAnswer, handler.userStore)).Methods(http.MethodPost)
}

func (handler *Handler) handleCreateQuestionAnswer(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.CreateQuestionAnswerPayload
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

	questionAnswerID, err := handler.questionAnswerStore.CreateQuestionAnswer(userID, payload.CapsuleID, payload.Prompt, payload.Answer)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]uint{"questionAnswerId": questionAnswerID})
}

func (handler *Handler) handleDeleteQuestionAnswer(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.DeleteQuestionAnswerPayload
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

	err = handler.questionAnswerStore.DeleteQuestionAnswer(userID, payload.CapsuleID, payload.QuestionAnswerID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, nil)
}
