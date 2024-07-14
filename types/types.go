package types

import (
	"mime/multipart"
	"time"
)

// ====================================================================
// User
// ====================================================================

type User struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Password      string    `json:"-"`
	ReferralCount uint      `json:"referralCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(userId uint) (*User, error)
	GetUserNameById(userId uint) (string, error)
	CreateUser(name string, email string, phone string, password string) error
	DeleteUser(userId uint) error
	UpdateUser(userId uint, name string, email string, phone string) error
	UpdateUserPassword(userId uint, password string) error
	ProcessContacts([]Contact) ([]Contact, []Contact, []Contact, error)
	AddReferral(userId uint, phone string) error
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,min=10,max=10"`
	Password string `json:"password" validate:"required,min=6,max=130"`
}

type UpdateUserPayload struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,min=10,max=10"`
}

type UpdateUserPasswordPayload struct {
	Password string `json:"password" validate:"required,min=6,max=130"`
}

type ProcessContactsPayload struct {
	Contacts []Contact `json:"contacts" validate:"required,dive"`
}

type AddReferralPayload struct {
	Phone string `json:"phone" validate:"required,min=10,max=10"`
}

type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Referral struct {
	ID            uint      `json:"id"`
	Phone         string    `json:"phone"`
	ReferralCount uint      `json:"referralCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

// ====================================================================
// File
// ====================================================================

type FileStore interface {
	UploadFile(userId uint, file multipart.File, fileHeader *multipart.FileHeader) (string, string, error)
	DeleteFile(objectName string) error
}

type DeleteFilePayload struct {
	ObjectName string `json:"objectName" validate:"required"`
}

// ====================================================================
// Capsule
// ====================================================================

type Capsule struct {
	ID             uint      `json:"id"`
	Code           string    `json:"code"`
	CreatedAt      time.Time `json:"createdAt"`
	Public         bool      `json:"public"`
	CapsuleOwnerID uint      `json:"capsuleOwnerId"`

	CapsuleMember1ID uint `json:"capsuleMember1Id"`
	CapsuleMember2ID uint `json:"capsuleMember2Id"`
	CapsuleMember3ID uint `json:"capsuleMember3Id"`
	CapsuleMember4ID uint `json:"capsuleMember4Id"`
	CapsuleMember5ID uint `json:"capsuleMember5Id"`

	CapsuleMember1Sealed bool `json:"capsuleMember1Sealed"`
	CapsuleMember2Sealed bool `json:"capsuleMember2Sealed"`
	CapsuleMember3Sealed bool `json:"capsuleMember3Sealed"`
	CapsuleMember4Sealed bool `json:"capsuleMember4Sealed"`
	CapsuleMember5Sealed bool `json:"capsuleMember5Sealed"`

	Vessel     string     `json:"vessel"`
	Name       string     `json:"name"`
	DateToOpen *time.Time `json:"dateToOpen"`
	EmailSent  bool       `json:"emailSent"`
	Sealed     string     `json:"sealed"`
}

type CapsuleStore interface {
	GetCapsules(userId uint) ([]Capsule, error)
	GetCapsuleById(userId uint, capsuleId uint) (Capsule, error)
	GetCapsuleByIdUnsafe(userId uint, capsuleId uint) (Capsule, error)
	CreateCapsule(userId uint, vessel string, public bool) (uint, error)
	JoinCapsule(userId uint, code string) error
	DeleteCapsule(userId uint, capsuleId uint) ([]string, error)
	NameCapsule(userId uint, capsuleId uint, name string) error
	SealCapsule(userId uint, capsuleId uint, dateToOpen time.Time) error
	MemberSealCapsule(userId uint, capsuleId uint, memberNumber uint) error
	OpenCapsule(userId uint, capsuleId uint) error
	SendReminderMail() error
}

type GetCapsuleByIdResponse struct {
	Capsule         Capsule          `json:"capsule"`
	Songs           []Song           `json:"songs"`
	QuestionAnswers []QuestionAnswer `json:"questionAnswers"`
	Writings        []Writing        `json:"writings"`
	Photos          []Photo          `json:"photos"`
	Audios          []Audio          `json:"audios"`
	Doodles         []Doodle         `json:"doodles"`
	MiscFiles       []MiscFile       `json:"miscFiles"`
}

type CreateCapsulePayload struct {
	Vessel string `json:"vessel" validate:"required,min=1,max=32"`
	Public bool   `json:"public"`
}

type JoinCapsulePayload struct {
	Code string `json:"code" validate:"required,min=10,max=10"`
}

type DeleteCapsulePayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
}

type NameCapsulePayload struct {
	CapsuleID uint   `json:"capsuleId" validate:"required"`
	Name      string `json:"name" validate:"required,min=1,max=255"`
}

type SealCapsulePayload struct {
	CapsuleID  uint   `json:"capsuleId" validate:"required"`
	DateToOpen string `json:"dateToOpen" validate:"required"`
}

type MemberSealCapsulePayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
}

type OpenCapsulePayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
}

// ====================================================================
// Song
// ====================================================================

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
	DeleteSong(userID uint, capsuleID uint, songID uint) error
}

type CreateSongPayload struct {
	CapsuleID   uint   `json:"capsuleId" validate:"required"`
	SpotifyID   string `json:"spotifyId" validate:"required"`
	Name        string `json:"name" validate:"required"`
	ArtistName  string `json:"artistName" validate:"required"`
	AlbumArtURL string `json:"albumArtURL"`
}

type DeleteSongPayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
	SongID    uint `json:"songId" validate:"required"`
}

// ====================================================================
// QuestionAnswer
// ====================================================================

type QuestionAnswer struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	CapsuleID uint      `json:"capsuleId"`
	Prompt    string    `json:"prompt"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"createdAt"`
}

type QuestionAnswerStore interface {
	GetQuestionAnswers(capsuleID uint) ([]QuestionAnswer, error)
	CreateQuestionAnswer(userID, capsuleID uint, prompt string, answer string) (uint, error)
	UpdateQuestionAnswer(userID, capsuleID uint, questionAnswerID uint, prompt string, answer string) error
	DeleteQuestionAnswer(userID uint, capsuleID uint, questionAnswerID uint) error
}

type CreateQuestionAnswerPayload struct {
	CapsuleID uint   `json:"capsuleId" validate:"required"`
	Prompt    string `json:"prompt" validate:"required,max=255"`
	Answer    string `json:"answer" validate:"required,max=1000"`
}

type UpdateQuestionAnswerPayload struct {
	QuestionAnswerID uint   `json:"questionAnswerId" validate:"required"`
	CapsuleID        uint   `json:"capsuleId" validate:"required"`
	Prompt           string `json:"prompt" validate:"required,max=255"`
	Answer           string `json:"answer" validate:"required,max=1000"`
}

type DeleteQuestionAnswerPayload struct {
	QuestionAnswerID uint `json:"questionAnswerId" validate:"required"`
	CapsuleID        uint `json:"capsuleId" validate:"required"`
}

// ====================================================================
// Writing
// ====================================================================

type Writing struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"userId"`
	CapsuleID uint      `json:"capsuleId"`
	Writing   string    `json:"writing"`
	CreatedAt time.Time `json:"createdAt"`
}

type WritingStore interface {
	GetWritings(capsuleID uint) ([]Writing, error)
	CreateWriting(userID, capsuleID uint, writing string) (uint, error)
	UpdateWriting(userID, capsuleID, writingID uint, writing string) error
	DeleteWriting(userID uint, capsuleID uint, writingID uint) error
}

type CreateWritingPayload struct {
	CapsuleID uint   `json:"capsuleId" validate:"required"`
	Writing   string `json:"writing" validate:"required,max=1000"`
}

type UpdateWritingPayload struct {
	CapsuleID uint   `json:"capsuleId" validate:"required"`
	WritingID uint   `json:"writingId" validate:"required"`
	Writing   string `json:"writing" validate:"required,max=1000"`
}

type DeleteWritingPayload struct {
	WritingID uint `json:"writingId" validate:"required"`
	CapsuleID uint `json:"capsuleId" validate:"required"`
}

// ====================================================================
// Photo
// ====================================================================

type Photo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	CapsuleID  uint      `json:"capsuleId"`
	ObjectName string    `json:"objectName"`
	FileURL    string    `json:"fileURL"`
	CreatedAt  time.Time `json:"createdAt"`
}

type PhotoStore interface {
	GetPhotos(capsuleID uint) ([]Photo, error)
	CreatePhoto(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error)
	DeletePhoto(userID uint, capsuleID uint, photoID uint) (string, error)
}

type CreatePhotoPayload struct {
	CapsuleID  uint   `json:"capsuleId" validate:"required"`
	ObjectName string `json:"objectName" validate:"required"`
	FileURL    string `json:"fileURL" validate:"required"`
}

type DeletePhotoPayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
	PhotoID   uint `json:"photoId" validate:"required"`
}

// ====================================================================
// Audio
// ====================================================================

type Audio struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	CapsuleID  uint      `json:"capsuleId"`
	ObjectName string    `json:"objectName"`
	FileURL    string    `json:"fileURL"`
	CreatedAt  time.Time `json:"createdAt"`
}

type AudioStore interface {
	GetAudios(capsuleID uint) ([]Audio, error)
	CreateAudio(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error)
	DeleteAudio(userID uint, capsuleID uint, audioID uint) (string, error)
}

type CreateAudioPayload struct {
	CapsuleID  uint   `json:"capsuleId" validate:"required"`
	ObjectName string `json:"objectName" validate:"required"`
	FileURL    string `json:"fileURL" validate:"required"`
}

type DeleteAudioPayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
	AudioID   uint `json:"audioId" validate:"required"`
}

// ====================================================================
// Doodle
// ====================================================================

type Doodle struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	CapsuleID  uint      `json:"capsuleId"`
	ObjectName string    `json:"objectName"`
	FileURL    string    `json:"fileURL"`
	CreatedAt  time.Time `json:"createdAt"`
}

type DoodleStore interface {
	GetDoodles(capsuleID uint) ([]Doodle, error)
	CreateDoodle(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error)
	DeleteDoodle(userID uint, capsuleID uint, doodleID uint) (string, error)
}

type CreateDoodlePayload struct {
	CapsuleID  uint   `json:"capsuleId" validate:"required"`
	ObjectName string `json:"objectName" validate:"required"`
	FileURL    string `json:"fileURL" validate:"required"`
}

type DeleteDoodlePayload struct {
	CapsuleID uint `json:"capsuleId" validate:"required"`
	DoodleID  uint `json:"doodleId" validate:"required"`
}

// ====================================================================
// MiscFile
// ====================================================================

type MiscFile struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	CapsuleID  uint      `json:"capsuleId"`
	ObjectName string    `json:"objectName"`
	FileURL    string    `json:"fileURL"`
	CreatedAt  time.Time `json:"createdAt"`
}

type MiscFileStore interface {
	GetMiscFiles(capsuleID uint) ([]MiscFile, error)
	CreateMiscFile(userID uint, capsuleID uint, objectName string, fileURL string) (uint, error)
	DeleteMiscFile(userID uint, capsuleID uint, miscFileID uint) (string, error)
}

type CreateMiscFilePayload struct {
	CapsuleID  uint   `json:"capsuleId" validate:"required"`
	ObjectName string `json:"objectName" validate:"required"`
	FileURL    string `json:"fileURL" validate:"required"`
}

type DeleteMiscFilePayload struct {
	CapsuleID  uint `json:"capsuleId" validate:"required"`
	MiscFileID uint `json:"miscFileId" validate:"required"`
}
