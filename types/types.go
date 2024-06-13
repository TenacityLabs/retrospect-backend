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
	CapsuleMember1ID *uint      `json:"capsuleMember1Id"`
	CapsuleMember2ID *uint      `json:"capsuleMember2Id"`
	CapsuleMember3ID *uint      `json:"capsuleMember3Id"`
	Vessel           *string    `json:"vessel"`
	DateToOpen       *time.Time `json:"dateToOpen"`
	EmailSent        bool       `json:"emailSent"`
	Sealed           bool       `json:"sealed"`
}

type CapsuleStore interface {
	GetCapsules(capsuleOwnerId uint) ([]Capsule, error)
	GetCapsuleById(capsuleOwnerId uint, capsuleId uint) (*Capsule, error)
}

type CreateCapsulePayload struct {
}
