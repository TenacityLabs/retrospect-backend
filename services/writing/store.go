package writing

import (
	"database/sql"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type WritingStore struct {
	db *sql.DB
}

func NewWritingStore(db *sql.DB) *WritingStore {
	return &WritingStore{
		db: db,
	}
}

func scanRowIntoWriting(row *sql.Rows) (*types.Writing, error) {
	writing := new(types.Writing)

	err := row.Scan(
		&writing.ID,
		&writing.UserID,
		&writing.CapsuleID,
		&writing.Writing,
		&writing.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return writing, nil
}

func (writingStore *WritingStore) GetWritings(capsuleID uint) ([]types.Writing, error) {
	rows, err := writingStore.db.Query("SELECT * FROM writings WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	writings := make([]types.Writing, 0)
	for rows.Next() {
		writing, err := scanRowIntoWriting(rows)
		if err != nil {
			return nil, err
		}
		writings = append(writings, *writing)
	}

	return writings, nil
}

func (writingStore *WritingStore) UpdateWriting(userID, capsuleID uint, writingID uint, writing string) error {
	_, err := writingStore.db.Exec("UPDATE writings SET writing = ? WHERE id = ? AND userId = ? AND capsuleId = ?", writing, writingID, userID, capsuleID)
	return err
}

func (writingStore *WritingStore) CreateWriting(userID uint, capsuleID uint, writing string) (uint, error) {
	res, err := writingStore.db.Exec("INSERT INTO writings (userId, capsuleId, writing) VALUES (?, ?, ?)", userID, capsuleID, writing)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (writingStore *WritingStore) DeleteWriting(userID uint, capsuleID uint, writingID uint) error {
	_, err := writingStore.db.Exec("DELETE FROM writings WHERE id = ? AND userId = ? AND capsuleId = ?", writingID, userID, capsuleID)
	return err
}
