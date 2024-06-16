package audio

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type AudioStore struct {
	db *sql.DB
}

func NewAudioStore(db *sql.DB) *AudioStore {
	return &AudioStore{
		db: db,
	}
}

func scanRowIntoAudio(row *sql.Rows) (*types.Audio, error) {
	audio := new(types.Audio)

	err := row.Scan(
		&audio.ID,
		&audio.UserID,
		&audio.CapsuleID,
		&audio.ObjectName,
		&audio.FileURL,
		&audio.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return audio, nil
}

func (audioStore *AudioStore) GetAudios(capsuleID uint) ([]types.Audio, error) {
	rows, err := audioStore.db.Query("SELECT * FROM audios WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	audios := make([]types.Audio, 0)
	for rows.Next() {
		audio, err := scanRowIntoAudio(rows)
		if err != nil {
			return nil, err
		}
		audios = append(audios, *audio)
	}

	return audios, nil
}

func (audioStore *AudioStore) CreateAudio(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error) {
	// check if audio already exists in capsule
	rows, err := audioStore.db.Query("SELECT * FROM audios WHERE userId = ? AND capsuleId = ? AND fileURL = ?", userID, capsuleID, fileURL)
	if err != nil {
		return 0, err
	}
	audio := new(types.Audio)
	for rows.Next() {
		audio, err = scanRowIntoAudio(rows)
		if err != nil {
			return 0, err
		}
	}
	if audio.FileURL == fileURL {
		return 0, fmt.Errorf("you already added this audio to the capsule")
	}

	res, err := audioStore.db.Exec("INSERT INTO audios (userId, capsuleId, objectName, fileURL) VALUES (?, ?, ?, ?)", userID, capsuleID, objectName, fileURL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (audioStore *AudioStore) DeleteAudio(userID uint, capsuleID uint, audioID uint) (string, error) {
	// find original audio
	rows, err := audioStore.db.Query("SELECT * FROM audios WHERE id = ? AND userId = ? AND capsuleId = ?", audioID, userID, capsuleID)
	if err != nil {
		return "", err
	}
	audio := new(types.Audio)
	for rows.Next() {
		audio, err = scanRowIntoAudio(rows)
		if err != nil {
			return "", err
		}
	}
	if audio.ID != audioID {
		return "", fmt.Errorf("audio not found")
	}

	_, err = audioStore.db.Exec("DELETE FROM audios WHERE id = ? AND userId = ? AND capsuleId = ?", audioID, userID, capsuleID)
	return audio.ObjectName, err
}
