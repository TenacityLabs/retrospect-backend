package user

import (
	"database/sql"
	"fmt"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func scanRowIntoUser(row *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userStore *UserStore) GetUserByEmail(email string) (*types.User, error) {
	rows, err := userStore.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (userStore *UserStore) GetUserById(userId uint) (*types.User, error) {
	rows, err := userStore.db.Query("SELECT * FROM users WHERE id = ?", userId)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID != userId {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (userStore *UserStore) CreateUser(user types.User) error {
	_, err := userStore.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
