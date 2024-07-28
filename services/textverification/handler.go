package textverification

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TenacityLabs/retrospect-backend/types"
	"github.com/TenacityLabs/retrospect-backend/utils"
	"github.com/gorilla/mux"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type Handler struct {
	store            *TextVerificationStore
	twilioClient     *twilio.RestClient
	verifyServiceSID string
}


func NewHandler(store *TextVerificationStore, twilioClient *twilio.RestClient, verifyServiceSID string) *Handler {
	return &Handler{
		store:            store,
		twilioClient:     twilioClient,
		verifyServiceSID: verifyServiceSID,
	}
}

func (h *Handler) HandleSendVerification(w http.ResponseWriter, r *http.Request) {
	
	var payload types.SendVerificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	log.Printf("HandleSendVerification was hit with payload: %+v", payload)

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	params := &verify.CreateVerificationParams{}
	params.SetTo(payload.Phone)
	params.SetChannel("sms")

	resp, err := h.twilioClient.VerifyV2.CreateVerification(h.verifyServiceSID, params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send verification: %v", err))
		return
	}

	if resp.Sid == nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send verification: no SID returned"))
		return
	}

	err = h.store.SendVerification(payload.Phone)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to send verification: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Verification code sent"})
}

func (h *Handler) HandleVerifyCode(w http.ResponseWriter, r *http.Request) {
	var payload types.VerifyCodePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(payload.Phone)
	params.SetCode(payload.Code)

	resp, err := h.twilioClient.VerifyV2.CreateVerificationCheck(h.verifyServiceSID, params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to verify code: %v", err))
		return
	}

	if resp.Status == nil || *resp.Status != "approved" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired verification code"))
		return
	}

	isValid, err := h.store.CheckVerification(payload.Phone, payload.Code)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to verify code: %v", err))
		return
	}

	if !isValid {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid or expired verification code"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Verification successful"})
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/text-verification/send", h.HandleSendVerification).Methods(http.MethodPost)
	router.HandleFunc("/text-verification/verify", h.HandleVerifyCode).Methods(http.MethodPost)
}