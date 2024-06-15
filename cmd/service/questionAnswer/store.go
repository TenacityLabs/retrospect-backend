package questionAnswer

import (
	"database/sql"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type QuestionAnswerStore struct {
	db *sql.DB
}

func NewQuestionAnswerStore(db *sql.DB) *QuestionAnswerStore {
	return &QuestionAnswerStore{
		db: db,
	}
}

func scanRowIntoQuestionAnswer(row *sql.Rows) (*types.QuestionAnswer, error) {
	questionAnswer := new(types.QuestionAnswer)

	err := row.Scan(
		&questionAnswer.ID,
		&questionAnswer.UserID,
		&questionAnswer.CapsuleID,
		&questionAnswer.Prompt,
		&questionAnswer.Answer,
		&questionAnswer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return questionAnswer, nil
}

func (questionAnswerStore *QuestionAnswerStore) GetQuestionAnswers(capsuleID uint) ([]types.QuestionAnswer, error) {
	rows, err := questionAnswerStore.db.Query("SELECT * FROM questionAnswers WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	questionAnswers := make([]types.QuestionAnswer, 0)
	for rows.Next() {
		questionAnswer, err := scanRowIntoQuestionAnswer(rows)
		if err != nil {
			return nil, err
		}
		questionAnswers = append(questionAnswers, *questionAnswer)
	}

	return questionAnswers, nil
}

func (questionAnswerStore *QuestionAnswerStore) CreateQuestionAnswer(userID uint, capsuleID uint, prompt string, answer string) (uint, error) {
	res, err := questionAnswerStore.db.Exec("INSERT INTO questionAnswers (userId, capsuleId, prompt, answer) VALUES (?, ?, ?, ?)", userID, capsuleID, prompt, answer)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (questionAnswerStore *QuestionAnswerStore) DeleteQuestionAnswer(userID uint, capsuleID uint, questionAnswerID uint) error {
	_, err := questionAnswerStore.db.Exec("DELETE FROM questionAnswers WHERE id = ? AND userId = ? AND capsuleId = ?", questionAnswerID, userID, capsuleID)
	return err
}
