package textverification

import (
	"database/sql"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type TextVerificationStore struct {
	db            *sql.DB
	twilioClient  *twilio.RestClient
	verifyService string
}

func NewTextVerificationStore(db *sql.DB, twilioClient *twilio.RestClient, verifyService string) *TextVerificationStore {
	return &TextVerificationStore{
		db:            db,
		twilioClient:  twilioClient,
		verifyService: verifyService,
	}
}

func (s *TextVerificationStore) SendVerification(phone string) error {
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")

	_, err := s.twilioClient.VerifyV2.CreateVerification(s.verifyService, params)
	return err
}

func (s *TextVerificationStore) CheckVerification(phone, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phone)
	params.SetCode(code)

	resp, err := s.twilioClient.VerifyV2.CreateVerificationCheck(s.verifyService, params)
	if err != nil {
		return false, err
	}

	return *resp.Status == "approved", nil
}