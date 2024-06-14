package types

import "time"

type User struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(userId uint) (*User, error)
	CreateUser(user User) error
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=130"`
}

type Capsule struct {
	ID               uint       `json:"id"`
	Code             string     `json:"code"`
	CreatedAt        time.Time  `json:"createdAt"`
	Public           bool       `json:"public"`
	CapsuleOwnerID   uint       `json:"capsuleOwnerId"`
	CapsuleMember1ID uint       `json:"capsuleMember1Id"`
	CapsuleMember2ID uint       `json:"capsuleMember2Id"`
	CapsuleMember3ID uint       `json:"capsuleMember3Id"`
	Vessel           string     `json:"vessel"`
	DateToOpen       *time.Time `json:"dateToOpen"`
	EmailSent        bool       `json:"emailSent"`
	Sealed           bool       `json:"sealed"`
}

type CapsuleStore interface {
	GetCapsules(userId uint) ([]Capsule, error)
	GetCapsuleById(userId uint, capsuleID uint) (*Capsule, error)
	CreateCapsule(userId uint, vessel string, public bool) (uint, error)
}

type CreateCapsulePayload struct {
	Vessel string `json:"vessel" valiedate:"required,min=1,max=32"`
	Public bool   `json:"public"`
}

type Song struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"userId"`
	CapsuleID   uint      `json:"capsuleId"`
	SpotifyID   string    `json:"spotifyId"`
	Name        string    `json:"name"`
	ArtistName  string    `json:"artistName"`
	AlbumArtURL string    `json:"albumArtURL"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SongStore interface {
	GetSongs(capsuleID uint) ([]Song, error)
	CreateSong(userID uint, capsuleID uint, spotifyID string, name string, artistName string, albumArtURL string) (uint, error)
}

type CreateSongPayload struct {
	CapsuleID   uint   `json:"capsuleId" validate:"required"`
	SpotifyID   string `json:"spotifyId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	ArtistName  string `json:"artistName" validate:"required"`
	AlbumArtURL string `json:"albumArtURL"`
}