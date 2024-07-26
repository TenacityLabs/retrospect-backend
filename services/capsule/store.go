package capsule

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/TenacityLabs/retrospect-backend/config"
	"github.com/TenacityLabs/retrospect-backend/types"
)

type CapsuleStore struct {
	db  *sql.DB
	rng *rand.Rand
}

func NewCapsuleStore(db *sql.DB) *CapsuleStore {
	return &CapsuleStore{
		db:  db,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func scanRowIntoCapsule(row *sql.Rows) (*types.Capsule, error) {
	capsule := new(types.Capsule)

	var capsuleMember1Id, capsuleMember2Id, capsuleMember3Id, capsuleMember4Id, capsuleMember5Id sql.NullInt64

	err := row.Scan(
		&capsule.ID,
		&capsule.Code,
		&capsule.CreatedAt,
		&capsule.Public,
		&capsule.CapsuleOwnerID,
		&capsuleMember1Id,
		&capsuleMember2Id,
		&capsuleMember3Id,
		&capsuleMember4Id,
		&capsuleMember5Id,
		&capsule.CapsuleMember1Sealed,
		&capsule.CapsuleMember2Sealed,
		&capsule.CapsuleMember3Sealed,
		&capsule.CapsuleMember4Sealed,
		&capsule.CapsuleMember5Sealed,
		&capsule.Vessel,
		&capsule.Name,
		&capsule.DateToOpen,
		&capsule.EmailSent,
		&capsule.Sealed,
	)
	if err != nil {
		return nil, err
	}

	// Set CapsuleMemberID fields to 0 if they are NULL
	if capsuleMember1Id.Valid {
		capsule.CapsuleMember1ID = uint(capsuleMember1Id.Int64)
	} else {
		capsule.CapsuleMember1ID = 0
	}
	if capsuleMember2Id.Valid {
		capsule.CapsuleMember2ID = uint(capsuleMember2Id.Int64)
	} else {
		capsule.CapsuleMember2ID = 0
	}
	if capsuleMember3Id.Valid {
		capsule.CapsuleMember3ID = uint(capsuleMember3Id.Int64)
	} else {
		capsule.CapsuleMember3ID = 0
	}
	if capsuleMember4Id.Valid {
		capsule.CapsuleMember4ID = uint(capsuleMember4Id.Int64)
	} else {
		capsule.CapsuleMember4ID = 0
	}
	if capsuleMember5Id.Valid {
		capsule.CapsuleMember5ID = uint(capsuleMember5Id.Int64)
	} else {
		capsule.CapsuleMember5ID = 0
	}

	return capsule, nil
}

func (capsuleStore *CapsuleStore) GetCapsules(userId uint) ([]types.Capsule, error) {
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE capsuleOwnerId = ? OR capsuleMember1Id = ? OR capsuleMember2Id = ? OR capsuleMember3Id = ? OR capsuleMember4Id = ? OR capsuleMember5Id = ?", userId, userId, userId, userId, userId, userId)
	if err != nil {
		return nil, err
	}

	capsules := make([]types.Capsule, 0)
	for rows.Next() {
		capsule, err := scanRowIntoCapsule(rows)
		if err != nil {
			return nil, err
		}
		capsules = append(capsules, *capsule)
	}

	return capsules, nil
}

func (capsuleStore *CapsuleStore) GetCapsuleById(userId uint, capsuleId uint) (types.Capsule, error) {
	capsule := new(types.Capsule)
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE id = ?", capsuleId)
	if err != nil {
		return *capsule, err
	}

	for rows.Next() {
		capsule, err = scanRowIntoCapsule(rows)
		if err != nil {
			return *capsule, err
		}
	}

	if capsule.ID != capsuleId {
		return *capsule, fmt.Errorf("capsule not found")
	}
	if capsule.CapsuleOwnerID != userId && capsule.CapsuleMember1ID != userId && capsule.CapsuleMember2ID != userId && capsule.CapsuleMember3ID != userId && capsule.CapsuleMember4ID != userId && capsule.CapsuleMember5ID != userId {
		return *capsule, fmt.Errorf("user is not authorized to view this capsule")
	}
	if capsule.Sealed == "sealed" || capsule.Sealed == "opened" {
		return *capsule, fmt.Errorf("capsule cannot be modified because it has already been sealed or opened")
	}

	return *capsule, nil
}

func (capsuleStore *CapsuleStore) GetCapsuleByIdUnsafe(userId uint, capsuleId uint) (types.Capsule, error) {
	capsule := new(types.Capsule)
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE id = ?", capsuleId)
	if err != nil {
		return *capsule, err
	}

	for rows.Next() {
		capsule, err = scanRowIntoCapsule(rows)
		if err != nil {
			return *capsule, err
		}
	}

	if capsule.ID != capsuleId {
		return *capsule, fmt.Errorf("capsule not found")
	}
	if capsule.CapsuleOwnerID != userId && capsule.CapsuleMember1ID != userId && capsule.CapsuleMember2ID != userId && capsule.CapsuleMember3ID != userId && capsule.CapsuleMember4ID != userId && capsule.CapsuleMember5ID != userId {
		return *capsule, fmt.Errorf("user is not authorized to view this capsule")
	}

	return *capsule, nil
}

func (capsuleStore *CapsuleStore) GenerateCapsuleCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[capsuleStore.rng.Intn(len(charset))]
	}
	return string(b)
}

func (capsuleStore *CapsuleStore) CreateCapsule(userId uint, vessel string, public bool) (uint, error) {
	// generate unique capulse code
	var code string
	generateCodeAttempts := 0
	const codeLength = 10

	for {
		code = capsuleStore.GenerateCapsuleCode(codeLength)
		var count int
		err := capsuleStore.db.QueryRow("SELECT COUNT(*) FROM capsules WHERE code = ?", code).Scan(&count)
		if err != nil {
			return 0, err
		}
		if count == 0 {
			break
		}
		if generateCodeAttempts > 10 {
			return 0, fmt.Errorf("failed to generate unique capsule code after 10 attempts")
		}
		generateCodeAttempts++
	}

	// check if vessel is valid
	allowedVessels := []string{"box", "suitcase", "guitar case", "bottle", "shoe", "garbage"}
	validVessel := false
	for _, allowedVessel := range allowedVessels {
		if vessel == allowedVessel {
			validVessel = true
			break
		}
	}
	if !validVessel {
		return 0, fmt.Errorf("invalid vessel")
	}

	res, err := capsuleStore.db.Exec("INSERT INTO capsules (code, capsuleOwnerId, vessel, name, public) VALUES (?, ?, ?, 'My Time Capsule', ?)", code, userId, vessel, public)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (capsuleStore *CapsuleStore) JoinCapsule(userId uint, code string) error {
	// get the capsule
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE code = ?", code)
	if err != nil {
		return err
	}

	capsule := new(types.Capsule)
	for rows.Next() {
		capsule, err = scanRowIntoCapsule(rows)
		if err != nil {
			return err
		}
	}
	if capsule.Code != code {
		return fmt.Errorf("capsule not found")
	}

	// check if the user is already a member of the capsule
	if capsule.CapsuleOwnerID == userId || capsule.CapsuleMember1ID == userId || capsule.CapsuleMember2ID == userId || capsule.CapsuleMember3ID == userId || capsule.CapsuleMember4ID == userId || capsule.CapsuleMember5ID == userId {
		return fmt.Errorf("you are already a member of the capsule")
	}

	// check for the first available member slot
	if capsule.CapsuleMember1ID == 0 {
		_, err = capsuleStore.db.Exec("UPDATE capsules SET capsuleMember1Id = ? WHERE id = ?", userId, capsule.ID)
	} else if capsule.CapsuleMember2ID == 0 {
		_, err = capsuleStore.db.Exec("UPDATE capsules SET capsuleMember2Id = ? WHERE id = ?", userId, capsule.ID)
	} else if capsule.CapsuleMember3ID == 0 {
		_, err = capsuleStore.db.Exec("UPDATE capsules SET capsuleMember3Id = ? WHERE id = ?", userId, capsule.ID)
	} else if capsule.CapsuleMember4ID == 0 {
		_, err = capsuleStore.db.Exec("UPDATE capsules SET capsuleMember4Id = ? WHERE id = ?", userId, capsule.ID)
	} else if capsule.CapsuleMember5ID == 0 {
		_, err = capsuleStore.db.Exec("UPDATE capsules SET capsuleMember5Id = ? WHERE id = ?", userId, capsule.ID)
	} else {
		return fmt.Errorf("capsule already has the maximum number of members")
	}

	return err
}

func (capsuleStore *CapsuleStore) DeleteCapsule(userId uint, capsuleId uint) ([]string, error) {
	var objectNames []string

	_, err := capsuleStore.db.Exec("DELETE FROM songs WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM questionAnswers WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM writings WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	// get all photo objectNames
	rows, err := capsuleStore.db.Query("SELECT objectName FROM photos WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}
	defer rows.Close()

	for rows.Next() {
		var objectName string
		if err := rows.Scan(&objectName); err != nil {
			return objectNames, err
		}
		objectNames = append(objectNames, objectName)
	}
	if err := rows.Err(); err != nil {
		return objectNames, err
	}

	// get all audios
	rows, err = capsuleStore.db.Query("SELECT objectName FROM audios WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}
	defer rows.Close()

	for rows.Next() {
		var objectName string
		if err := rows.Scan(&objectName); err != nil {
			return objectNames, err
		}
		objectNames = append(objectNames, objectName)
	}
	if err := rows.Err(); err != nil {
		return objectNames, err
	}

	// get all doodles
	rows, err = capsuleStore.db.Query("SELECT objectName FROM doodles WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}
	defer rows.Close()

	for rows.Next() {
		var objectName string
		if err := rows.Scan(&objectName); err != nil {
			return objectNames, err
		}
		objectNames = append(objectNames, objectName)
	}
	if err := rows.Err(); err != nil {
		return objectNames, err
	}

	// get all misc files
	rows, err = capsuleStore.db.Query("SELECT objectName FROM miscFiles WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}
	defer rows.Close()

	for rows.Next() {
		var objectName string
		if err := rows.Scan(&objectName); err != nil {
			return objectNames, err
		}
		objectNames = append(objectNames, objectName)
	}
	if err := rows.Err(); err != nil {
		return objectNames, err
	}

	// delete all photos, audios, doodles, and misc files
	_, err = capsuleStore.db.Exec("DELETE FROM photos WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM audios WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM doodles WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM miscFiles WHERE capsuleId = ?", capsuleId)
	if err != nil {
		return objectNames, err
	}

	_, err = capsuleStore.db.Exec("DELETE FROM capsules WHERE id = ? AND capsuleOwnerId = ?", capsuleId, userId)
	return objectNames, err
}

func (capsuleStore *CapsuleStore) NameCapsule(userId uint, capsuleId uint, name string) error {
	_, err := capsuleStore.db.Exec("UPDATE capsules SET name = ? WHERE id = ? AND capsuleOwnerId = ?", name, capsuleId, userId)
	return err
}

func (capsuleStore *CapsuleStore) SealCapsule(userId uint, capsuleId uint, dateToOpen time.Time) error {
	_, err := capsuleStore.db.Exec("UPDATE capsules SET sealed = 'sealed', dateToOpen = ? WHERE id = ? AND capsuleOwnerId = ?", dateToOpen, capsuleId, userId)
	return err
}

func (capsuleStore *CapsuleStore) MemberSealCapsule(userId uint, capsuleId uint, memberNumber uint) error {
	if memberNumber < 1 || memberNumber > 5 {
		return fmt.Errorf("invalid member number")
	}

	query := fmt.Sprintf("UPDATE capsules SET capsuleMember%dSealed = TRUE WHERE id = ? AND capsuleMember%dId = ?", memberNumber, memberNumber)
	_, err := capsuleStore.db.Exec(query, capsuleId, userId)
	return err
}

func (capsuleStore *CapsuleStore) OpenCapsule(userId uint, capsuleId uint) error {
	_, err := capsuleStore.db.Exec("UPDATE capsules SET sealed = 'opened' WHERE id = ? AND capsuleOwnerId = ?", capsuleId, userId)
	return err
}

func (capsuleStore *CapsuleStore) SendReminderMail() error {
	findMailingListQuery := `
		SELECT c.id, u.email
		FROM capsules c
		JOIN users u ON c.capsuleOwnerId = u.id
		WHERE c.sealed = 'sealed' AND c.dateToOpen < NOW() AND c.emailSent = FALSE
		LIMIT 490
	`

	rows, err := capsuleStore.db.Query(findMailingListQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	emails := make([]string, 0)
	uniqueEmails := make(map[string]bool)
	capsuleIds := make([]uint, 0)
	for rows.Next() {
		var email string
		var capsuleId uint
		if err := rows.Scan(&capsuleId, &email); err != nil {
			return err
		}
		if !uniqueEmails[email] {
			emails = append(emails, email)
			uniqueEmails[email] = true
		}
		capsuleIds = append(capsuleIds, capsuleId)
	}

	if len(emails) > 0 {
		auth := smtp.PlainAuth(
			"",
			"retrospect.space@gmail.com",
			config.Envs.GmailAppPassword,
			"smtp.gmail.com",
		)

		msg := "Subject: Your Time Capsule is Ready!\n\nYour time capsule is ready to be opened! Open our app to see what's inside!"

		err = smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"retrospect.space@gmail.com",
			emails,
			[]byte(msg),
		)
		if err != nil {
			return err
		}

		capsuleStringIds := make([]string, len(capsuleIds))
		for i, capsuleId := range capsuleIds {
			capsuleStringIds[i] = fmt.Sprint(capsuleId)
		}
		// update db to mark capsules as emailSent
		updateEmailSentQuery := `
			UPDATE capsules
			SET emailSent = TRUE
			WHERE id IN (` + strings.Join(capsuleStringIds, ",") + `)
		`
		_, err = capsuleStore.db.Exec(updateEmailSentQuery)
		if err != nil {
			return err
		}
	}

	return nil
}
