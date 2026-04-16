// File: internal/repository/kursi_repository.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nickyrolly/ichat/internal/database"
)

type KursiRepository struct {
	db *database.DB
}

func NewKursiRepository(db *database.DB) *KursiRepository {
	return &KursiRepository{db: db}
}

type Kursi struct {
	ID         uuid.UUID
	Baris      int
	NomorKursi int
	TribunID   uuid.UUID
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
}

func (r *KursiRepository) Create(kursi *Kursi) error {
	query := `INSERT INTO kursi (id, baris, nomor_kursi, tribun_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, kursi.ID, kursi.Baris, kursi.NomorKursi, kursi.TribunID)
	if err != nil {
		return fmt.Errorf("failed to create kursi: %w", err)
	}
	return nil
}

func (r *KursiRepository) GetAll() ([]Kursi, error) {
	query := `SELECT id, baris, nomor_kursi, tribun_id, created_at, updated_at FROM kursi ORDER BY tribun_id, baris, nomor_kursi`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get kursis: %w", err)
	}
	defer rows.Close()

	var kursis []Kursi
	for rows.Next() {
		var kursi Kursi
		err := rows.Scan(&kursi.ID, &kursi.Baris, &kursi.NomorKursi, &kursi.TribunID, &kursi.CreatedAt, &kursi.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan kursi: %w", err)
		}
		kursis = append(kursis, kursi)
	}

	return kursis, nil
}

func (r *KursiRepository) GetCount() (int, error) {
	query := `SELECT COUNT(*) FROM kursi`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get kursi count: %w", err)
	}
	return count, nil
}

func (r *KursiRepository) GetByTribunID(tribunID uuid.UUID) ([]Kursi, error) {
	query := `SELECT id, baris, nomor_kursi, tribun_id, created_at, updated_at FROM kursi WHERE tribun_id = $1 ORDER BY baris, nomor_kursi`
	rows, err := r.db.Query(query, tribunID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kursis by tribun ID: %w", err)
	}
	defer rows.Close()

	var kursis []Kursi
	for rows.Next() {
		var kursi Kursi
		err := rows.Scan(&kursi.ID, &kursi.Baris, &kursi.NomorKursi, &kursi.TribunID, &kursi.CreatedAt, &kursi.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan kursi: %w", err)
		}
		kursis = append(kursis, kursi)
	}

	return kursis, nil
}
