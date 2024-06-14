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

	err := row.Scan(
		&capsule.ID,
		&capsule.Code,
		&capsule.CreatedAt,
		&capsule.Public,
		&capsule.CapsuleOwnerID,
		&capsule.CapsuleMember1ID,
		&capsule.CapsuleMember2ID,
		&capsule.CapsuleMember3ID,
		&capsule.Vessel,
		&capsule.DateToOpen,
		&capsule.EmailSent,
		&capsule.Sealed,
	)
	if err != nil {
		return nil, err
	}

	return capsule, nil
}

func (capsuleStore *CapsuleStore) GetCapsules(capsuleOwnerId uint) ([]types.Capsule, error) {
	rows, err := capsuleStore.db.Query("SELECT * FROM capsules WHERE capsuleOwnerId = ?", capsuleOwnerId)
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

func (capsuleStore *CapsuleStore) GetCapsuleById(capsuleOwnerId uint, capsuleId uint) (*types.Capsule, error) {
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

	return capsule, nil
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"
const codeLength = 10

func (capsuleStore *CapsuleStore) GenerateCapsuleCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[capsuleStore.rng.Intn(len(charset))]
	}
	return string(b)
}

func (capsuleStore *CapsuleStore) CreateCapsule(capsuleOwnerID uint, vessel string, public bool) (uint, error) {
	// generate unique capulse code
	var capsuleCode string
	generateCodeAttempts := 0
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

	res, err := capsuleStore.db.Exec("INSERT INTO capsules (code, capsuleOwnerId, vessel, public) VALUES (?, ?, ?, ?)", capsuleCode, capsuleOwnerID, vessel, public)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
