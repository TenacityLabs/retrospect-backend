package song

import (
	"database/sql"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type SongStore struct {
	db *sql.DB
}

func NewSongStore(db *sql.DB) *SongStore {
	return &SongStore{
		db: db,
	}
}

func scanRowIntoSong(row *sql.Rows) (*types.Song, error) {
	song := new(types.Song)

	err := row.Scan(
		&song.ID,
		&song.UserID,
		&song.CapsuleID,
		&song.SpotifyID,
		&song.Name,
		&song.ArtistName,
		&song.AlbumArtURL,
	)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (songStore *SongStore) GetSongs(capsuleID uint) ([]types.Song, error) {
	rows, err := songStore.db.Query("SELECT * FROM songs WHERE capsuleId = ?", capsuleID)
	if err != nil {
		return nil, err
	}

	songs := make([]types.Song, 0)
	for rows.Next() {
		song, err := scanRowIntoSong(rows)
		if err != nil {
			return nil, err
		}
		songs = append(songs, *song)
	}

	return songs, nil
}

func (songStore *SongStore) CreateSong(userID uint, capsuleID uint, spotifyID string, name string, artistName string, albumArtURL string) (uint, error) {
	res, err := songStore.db.Exec("INSERT INTO songs (userId, capsuleId, spotifyId, name, artistName, albumArtURL) VALUES (?, ?, ?, ?, ?, ?)", userID, capsuleID, spotifyID, name, artistName, albumArtURL)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
