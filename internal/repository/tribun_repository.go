// File: internal/repository/tribun_repository.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickyrolly/ichat/internal/database"
)

type TribunRepository struct {
	db *database.DB
}

func NewTribunRepository(db *database.DB) *TribunRepository {
	return &TribunRepository{db: db}
}

type Tribun struct {
	ID        uuid.UUID
	Nama      string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (r *TribunRepository) Create(tribun *Tribun) error {
	query := `INSERT INTO tribun (id, nama) VALUES ($1, $2)`
	_, err := r.db.Exec(query, tribun.ID, tribun.Nama)
	if err != nil {
		return fmt.Errorf("failed to create tribun: %w", err)
	}
	return nil
}

func (r *TribunRepository) GetAll() ([]Tribun, error) {
	query := `SELECT id, nama, created_at, updated_at FROM tribun ORDER BY nama`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tribuns: %w", err)
	}
	defer rows.Close()

	var tribuns []Tribun
	for rows.Next() {
		var tribun Tribun
		err := rows.Scan(&tribun.ID, &tribun.Nama, &tribun.CreatedAt, &tribun.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tribun: %w", err)
		}
		tribuns = append(tribuns, tribun)
	}

	return tribuns, nil
}

func (r *TribunRepository) GetByID(id uuid.UUID) (*Tribun, error) {
	query := `SELECT id, nama, created_at, updated_at FROM tribun WHERE id = $1`
	var tribun Tribun
	err := r.db.QueryRow(query, id).Scan(&tribun.ID, &tribun.Nama, &tribun.CreatedAt, &tribun.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get tribun by ID: %w", err)
	}
	return &tribun, nil
}

func (r *TribunRepository) GetCount() (int, error) {
	query := `SELECT COUNT(*) FROM tribun`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get tribun count: %w", err)
	}
	return count, nil
}
