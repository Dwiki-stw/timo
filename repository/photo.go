package repository

import (
	"context"
	"database/sql"
	"timo/domain"
	"timo/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type photo struct {
	pool *pgxpool.Pool
}

func NewPhoto(pool *pgxpool.Pool) domain.PhotoRepository {
	return &photo{pool: pool}
}

func (p *photo) Create(ctx context.Context, photo *models.Photo) error {
	query := `
		INSERT INTO photos (journal_id, url)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return p.pool.QueryRow(ctx, query, photo.JournalID, photo.Url).
		Scan(&photo.ID, &photo.CreatedAt)
}

func (p *photo) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM photos
		WHERE id = $1
	`

	result, err := p.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *photo) GetByJournalID(ctx context.Context, journalID int64) ([]models.Photo, error) {
	var photos []models.Photo
	query := `
		SELECT id, journal_id, url, created_at
		FROM photos
		WHERE journal_id = $1
	`

	rows, err := p.pool.Query(ctx, query, journalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Photo
		err := rows.Scan(&p.ID, &p.JournalID, &p.Url, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		photos = append(photos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
