package capsule

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/TenacityLabs/time-capsule-backend/types"
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

	var capsuleMember1Id, capsuleMember2Id, capsuleMember3Id sql.NullInt64

	err := row.Scan(
		&capsule.ID,
		&capsule.Code,
		&capsule.CreatedAt,
		&capsule.Public,
		&capsule.CapsuleOwnerID,
		&capsuleMember1Id,
		&capsuleMember2Id,
		&capsuleMember3Id,
		&capsule.Vessel,
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

	return capsule, nil
}

func (capsuleStore *CapsuleStore) GetCapsules(userId uint) ([]types.Capsule, error) {
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE capsuleOwnerId = ? OR capsuleMember1Id = ? OR capsuleMember2Id = ? OR capsuleMember3Id = ?", userId, userId, userId, userId)
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

func (capsuleStore *CapsuleStore) GetCapsuleById(userId uint, capsuleId uint) (*types.Capsule, error) {
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE id = ?", capsuleId)
	if err != nil {
		return nil, err
	}

	capsule := new(types.Capsule)
	for rows.Next() {
		capsule, err = scanRowIntoCapsule(rows)
		if err != nil {
			return nil, err
		}
	}

	if capsule.ID != capsuleId {
		return nil, fmt.Errorf("capsule not found")
	}
	if capsule.CapsuleOwnerID != userId && capsule.CapsuleMember1ID != userId && capsule.CapsuleMember2ID != userId && capsule.CapsuleMember3ID != userId {
		return nil, fmt.Errorf("user is not authorized to view this capsule")
	}

	return capsule, nil
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
	var capsuleCode string
	generateCodeAttempts := 0
	const codeLength = 10

	for {
		capsuleCode = capsuleStore.GenerateCapsuleCode(codeLength)
		var count int
		err := capsuleStore.db.QueryRow("SELECT COUNT(*) FROM capsules WHERE code = ?", capsuleCode).Scan(&count)
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

	res, err := capsuleStore.db.Exec("INSERT INTO capsules (code, capsuleOwnerId, vessel, public) VALUES (?, ?, ?, ?)", capsuleCode, userId, vessel, public)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
