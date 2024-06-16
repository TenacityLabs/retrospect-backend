package photo

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type PhotoStore struct {
	db *sql.DB
}

func NewPhotoStore(db *sql.DB) *PhotoStore {
	return &PhotoStore{
		db: db,
	}
}

func scanRowIntoPhoto(row *sql.Rows) (*types.Photo, error) {
	photo := new(types.Photo)

	err := row.Scan(
		&photo.ID,
		&photo.UserID,
		&photo.CapsuleID,
		&photo.ObjectName,
		&photo.FileURL,
		&photo.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return photo, nil
}

func (photoStore *PhotoStore) GetPhotos(capsuleID uint) ([]types.Photo, error) {
	rows, err := photoStore.db.Query("SELECT * FROM photos WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	photos := make([]types.Photo, 0)
	for rows.Next() {
		photo, err := scanRowIntoPhoto(rows)
		if err != nil {
			return nil, err
		}
		photos = append(photos, *photo)
	}

	return photos, nil
}

func (photoStore *PhotoStore) CreatePhoto(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error) {
	// check if photo already exists in capsule
	rows, err := photoStore.db.Query("SELECT * FROM photos WHERE userId = ? AND capsuleId = ? AND fileURL = ?", userID, capsuleID, fileURL)
	if err != nil {
		return 0, err
	}
	photo := new(types.Photo)
	for rows.Next() {
		photo, err = scanRowIntoPhoto(rows)
		if err != nil {
			return 0, err
		}
	}
	if photo.FileURL == fileURL {
		return 0, fmt.Errorf("you already added this photo to the capsule")
	}

	res, err := photoStore.db.Exec("INSERT INTO photos (userId, capsuleId, objectName, fileURL) VALUES (?, ?, ?, ?)", userID, capsuleID, objectName, fileURL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (photoStore *PhotoStore) DeletePhoto(userID uint, capsuleID uint, photoID uint) (string, error) {
	// find original photo
	rows, err := photoStore.db.Query("SELECT * FROM photos WHERE id = ? AND userId = ? AND capsuleId = ?", photoID, userID, capsuleID)
	if err != nil {
		return "", err
	}
	photo := new(types.Photo)
	for rows.Next() {
		photo, err = scanRowIntoPhoto(rows)
		if err != nil {
			return "", err
		}
	}
	if photo.ID != photoID {
		return "", fmt.Errorf("photo not found")
	}

	_, err = photoStore.db.Exec("DELETE FROM photos WHERE id = ? AND userId = ? AND capsuleId = ?", photoID, userID, capsuleID)
	return photo.ObjectName, err
}
