package user

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/TenacityLabs/retrospect-backend/types"
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
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.ReferralCount,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func scanRowsIntoPhones(rows *sql.Rows) ([]string, error) {
	var phones []string
	for rows.Next() {
		var phone string
		err := rows.Scan(&phone)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return phones, nil
}

func scanRowsIntoReferrals(rows *sql.Rows) ([]types.Referral, error) {
	var referrals []types.Referral
	for rows.Next() {
		var referral types.Referral
		err := rows.Scan(
			&referral.ID,
			&referral.Phone,
			&referral.ReferralCount,
			&referral.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		referrals = append(referrals, referral)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return referrals, nil
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

func (userStore *UserStore) CreateUser(name string, email string, phone string, password string) error {
	_, err := userStore.db.Exec("INSERT INTO users (name, email, phone, password) VALUES (?, ?, ?, ?)", name, email, phone, password)
	if err != nil {
		return err
	}
	return nil
}

// delete user leaves memory leaks (eg. capsules where the owner is deleted)
// but this feature is only for testing, so it's fine
func (userStore *UserStore) DeleteUser(userId uint) error {
	_, err := userStore.db.Exec("DELETE FROM users WHERE id = ?", userId)
	if err != nil {
		return err
	}
	return nil
}

func (userStore *UserStore) UpdateUser(userId uint, name string, email string, phone string) error {
	_, err := userStore.db.Exec("UPDATE users SET name = ?, email = ?, phone = ? WHERE id = ?", name, email, phone, userId)
	if err != nil {
		return err
	}
	return nil
}

func (userStore *UserStore) UpdateUserPassword(userId uint, password string) error {
	_, err := userStore.db.Exec("UPDATE users SET password = ? WHERE id = ?", password, userId)
	if err != nil {
		return err
	}
	return nil
}

func generatePlaceholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?,", n-1) + "?"
}

func convertToInterfaceSlice(stringSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(stringSlice))
	for i, v := range stringSlice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func (userStore *UserStore) ProcessContacts(contacts []types.Contact) ([]types.Contact, []types.Contact, []types.Contact, error) {
	// first strip all contacts that are already users
	var allPhones []string
	for _, contact := range contacts {
		if contact.Phone != "" {
			allPhones = append(allPhones, contact.Phone)
		}
	}

	var existingPhones []string
	if len(allPhones) != 0 {
		query := fmt.Sprintf("SELECT phone FROM users WHERE phone IN (%s)", generatePlaceholders(len(allPhones)))
		rows, err := userStore.db.Query(query, convertToInterfaceSlice(allPhones)...)
		if err != nil {
			return nil, nil, nil, err
		}
		existingPhones, err = scanRowsIntoPhones(rows)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	var unregisteredPhones []string
	for _, phone := range allPhones {
		var exists bool = false
		for _, existingPhone := range existingPhones {
			if phone == existingPhone {
				exists = true
				break
			}
		}
		if !exists {
			unregisteredPhones = append(unregisteredPhones, phone)
		}
	}

	var matchedReferrals []types.Referral
	if len(unregisteredPhones) != 0 {
		// query := fmt.Sprintf("SELECT phone, referralCount FROM referrals WHERE phone IN (%s)", generatePlaceholders(len(unregisteredPhones)))
		// rows, err := userStore.db.Query(query, convertToInterfaceSlice(unregisteredPhones)...)
		query := "SELECT * FROM referrals"
		rows, err := userStore.db.Query(query)
		if err != nil {
			return nil, nil, nil, err
		}
		matchedReferrals, err = scanRowsIntoReferrals(rows)
		if err != nil {
			return nil, nil, nil, err
		}

		var count int
		err = userStore.db.QueryRow("SELECT COUNT(*) FROM referrals").Scan(&count)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to count referrals: %w", err)
		}
	}

	var priorityPhones []string
	for _, referral := range matchedReferrals {
		if referral.ReferralCount < 10 {
			priorityPhones = append(priorityPhones, referral.Phone)
		}
	}

	var freshPhones []string
	for _, phone := range unregisteredPhones {
		var exists bool = false
		for _, priorityPhone := range priorityPhones {
			if phone == priorityPhone {
				exists = true
				break
			}
		}
		if !exists {
			freshPhones = append(freshPhones, phone)
		}
	}

	// populate back to contacts
	var existingContacts []types.Contact = make([]types.Contact, 0)
	for _, phone := range existingPhones {
		// find the contact with the phone number
		for _, contact := range contacts {
			if contact.Phone == phone {
				existingContacts = append(existingContacts, contact)
				break
			}
		}
	}
	var priorityContacts []types.Contact = make([]types.Contact, 0)
	for _, phone := range priorityPhones {
		// find the contact with the phone number
		for _, contact := range contacts {
			if contact.Phone == phone {
				priorityContacts = append(priorityContacts, contact)
				break
			}
		}
	}
	var freshContacts []types.Contact = make([]types.Contact, 0)
	for _, phone := range freshPhones {
		// find the contact with the phone number
		for _, contact := range contacts {
			if contact.Phone == phone {
				freshContacts = append(freshContacts, contact)
				break
			}
		}
	}

	return existingContacts, priorityContacts, freshContacts, nil
}

func (userStore *UserStore) AddReferral(userId uint, phone string) error {
	query := `
        INSERT INTO referrals (phone, referralCount) 
        VALUES (?, 1) 
        ON DUPLICATE KEY UPDATE referralCount = referralCount + 1
    `
	_, err := userStore.db.Exec(query, phone)
	if err != nil {
		return err
	}

	_, err = userStore.db.Exec("UPDATE users SET referralCount = referralCount + 1 WHERE id = ?", userId)
	if err != nil {
		return err
	}

	return nil
}
