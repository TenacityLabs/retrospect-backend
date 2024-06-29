package miscFile

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/retrospect-backend/types"
)

type MiscFileStore struct {
	db *sql.DB
}

func NewMiscFileStore(db *sql.DB) *MiscFileStore {
	return &MiscFileStore{
		db: db,
	}
}

func scanRowIntoMiscFile(row *sql.Rows) (*types.MiscFile, error) {
	miscFile := new(types.MiscFile)

	err := row.Scan(
		&miscFile.ID,
		&miscFile.UserID,
		&miscFile.CapsuleID,
		&miscFile.ObjectName,
		&miscFile.FileURL,
		&miscFile.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return miscFile, nil
}

func (miscFileStore *MiscFileStore) GetMiscFiles(capsuleID uint) ([]types.MiscFile, error) {
	rows, err := miscFileStore.db.Query("SELECT * FROM miscFiles WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	miscFiles := make([]types.MiscFile, 0)
	for rows.Next() {
		miscFile, err := scanRowIntoMiscFile(rows)
		if err != nil {
			return nil, err
		}
		miscFiles = append(miscFiles, *miscFile)
	}

	return miscFiles, nil
}

func (miscFileStore *MiscFileStore) CreateMiscFile(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error) {
	// check if miscFile already exists in capsule
	rows, err := miscFileStore.db.Query("SELECT * FROM miscFiles WHERE userId = ? AND capsuleId = ? AND fileURL = ?", userID, capsuleID, fileURL)
	if err != nil {
		return 0, err
	}
	miscFile := new(types.MiscFile)
	for rows.Next() {
		miscFile, err = scanRowIntoMiscFile(rows)
		if err != nil {
			return 0, err
		}
	}
	if miscFile.FileURL == fileURL {
		return 0, fmt.Errorf("you already added this miscFile to the capsule")
	}

	res, err := miscFileStore.db.Exec("INSERT INTO miscFiles (userId, capsuleId, objectName, fileURL) VALUES (?, ?, ?, ?)", userID, capsuleID, objectName, fileURL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// TODO: edit miscFiles?

func (miscFileStore *MiscFileStore) DeleteMiscFile(userID uint, capsuleID uint, miscFileID uint) (string, error) {
	// find original miscFile
	rows, err := miscFileStore.db.Query("SELECT * FROM miscFiles WHERE id = ? AND userId = ? AND capsuleId = ?", miscFileID, userID, capsuleID)
	if err != nil {
		return "", err
	}
	miscFile := new(types.MiscFile)
	for rows.Next() {
		miscFile, err = scanRowIntoMiscFile(rows)
		if err != nil {
			return "", err
		}
	}
	if miscFile.ID != miscFileID {
		return "", fmt.Errorf("miscFile not found")
	}

	_, err = miscFileStore.db.Exec("DELETE FROM miscFiles WHERE id = ? AND userId = ? AND capsuleId = ?", miscFileID, userID, capsuleID)
	return miscFile.ObjectName, err
}
