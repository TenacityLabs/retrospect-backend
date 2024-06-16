package doodle

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type DoodleStore struct {
	db *sql.DB
}

func NewDoodleStore(db *sql.DB) *DoodleStore {
	return &DoodleStore{
		db: db,
	}
}

func scanRowIntoDoodle(row *sql.Rows) (*types.Doodle, error) {
	doodle := new(types.Doodle)

	err := row.Scan(
		&doodle.ID,
		&doodle.UserID,
		&doodle.CapsuleID,
		&doodle.ObjectName,
		&doodle.FileURL,
		&doodle.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return doodle, nil
}

func (doodleStore *DoodleStore) GetDoodles(capsuleID uint) ([]types.Doodle, error) {
	rows, err := doodleStore.db.Query("SELECT * FROM doodles WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	doodles := make([]types.Doodle, 0)
	for rows.Next() {
		doodle, err := scanRowIntoDoodle(rows)
		if err != nil {
			return nil, err
		}
		doodles = append(doodles, *doodle)
	}

	return doodles, nil
}

func (doodleStore *DoodleStore) CreateDoodle(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error) {
	// check if doodle already exists in capsule
	rows, err := doodleStore.db.Query("SELECT * FROM doodles WHERE userId = ? AND capsuleId = ? AND fileURL = ?", userID, capsuleID, fileURL)
	if err != nil {
		return 0, err
	}
	doodle := new(types.Doodle)
	for rows.Next() {
		doodle, err = scanRowIntoDoodle(rows)
		if err != nil {
			return 0, err
		}
	}
	if doodle.FileURL == fileURL {
		return 0, fmt.Errorf("you already added this doodle to the capsule")
	}

	res, err := doodleStore.db.Exec("INSERT INTO doodles (userId, capsuleId, objectName, fileURL) VALUES (?, ?, ?, ?)", userID, capsuleID, objectName, fileURL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// TODO: edit doodles?

func (doodleStore *DoodleStore) DeleteDoodle(userID uint, capsuleID uint, doodleID uint) (string, error) {
	// find original doodle
	rows, err := doodleStore.db.Query("SELECT * FROM doodles WHERE id = ? AND userId = ? AND capsuleId = ?", doodleID, userID, capsuleID)
	if err != nil {
		return "", err
	}
	doodle := new(types.Doodle)
	for rows.Next() {
		doodle, err = scanRowIntoDoodle(rows)
		if err != nil {
			return "", err
		}
	}
	if doodle.ID != doodleID {
		return "", fmt.Errorf("doodle not found")
	}

	_, err = doodleStore.db.Exec("DELETE FROM doodles WHERE id = ? AND userId = ? AND capsuleId = ?", doodleID, userID, capsuleID)
	return doodle.ObjectName, err
}
