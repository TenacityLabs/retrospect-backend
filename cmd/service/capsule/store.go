package capsule

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/TenacityLabs/time-capsule-backend/types"
)

type Store struct {
	db  *sql.DB
	rng *rand.Rand
}

func NewStore(db *sql.DB) *Store {
	return &Store{
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

func (store *Store) GetCapsules(capsuleOwnerId uint) ([]types.Capsule, error) {
	rows, err := store.db.Query("SELECT * FROM capsules WHERE capsuleOwnerId = ?", capsuleOwnerId)
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

func (store *Store) GetCapsuleById(capsuleOwnerId uint, capsuleId uint) (*types.Capsule, error) {
	rows, err := store.db.Query("SELECT * FROM capsules WHERE id = ?", capsuleId)
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

func (store *Store) generateCapsuleCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[store.rng.Intn(len(charset))]
	}
	return string(b)
}

func (store *Store) CreateCapsule(capsuleOwnerId uint) error {
	var capsuleCode string
	for {
		capsuleCode = store.generateCapsuleCode(codeLength)
		var count int
		err := store.db.QueryRow("SELECT COUNT(*) FROM capsules WHERE code = ?", capsuleCode).Scan(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			break
		}
	}

	_, err := store.db.Exec("INSERT INTO capsules (code, capsuleOwnerId) VALUES (?, ?)", capsuleCode, capsuleOwnerId)
	if err != nil {
		return err
	}
	return nil
}
