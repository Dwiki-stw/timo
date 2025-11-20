package repository

import (
	"context"
	"database/sql"
	"time"
	"timo/domain"
	"timo/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type journal struct {
	pool *pgxpool.Pool
}

func NewJournal(pool *pgxpool.Pool) domain.JournalRepository {
	return &journal{pool: pool}
}

func (j *journal) Create(ctx context.Context, journal *models.Journal) error {
	query := `
		INSERT INTO journals (user_id, title, text, mood_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, uid, created_at, updated_at
	`

	return j.pool.QueryRow(ctx, query, journal.UserID, journal.Title, journal.Text, journal.MoodID).
		Scan(&journal.ID, &journal.Uid, &journal.CreatedAt, &journal.UpdatedAt)
}

func (j *journal) Delete(ctx context.Context, uid string) error {
	query := `
		DELETE FROM journals WHERE uid = $1
	`

	result, err := j.pool.Exec(ctx, query, uid)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (j *journal) GetByID(ctx context.Context, uid string) (*models.Journal, error) {
	var journal models.Journal

	query := `
		SELECT j.id, j.uid, j.user_id, j.title, j.text, j.mood_id, m.label AS mood_label,  j.created_at, j.updated_at
		FROM journals j
		JOIN moods m ON m.id = j.mood_id
		WHERE uid = $1
	`

	err := j.pool.QueryRow(ctx, query, uid).
		Scan(&journal.ID, &journal.Uid, &journal.UserID, &journal.Title, &journal.Text, &journal.MoodID, &journal.MoodLabel, &journal.CreatedAt, &journal.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &journal, nil
}

func (j *journal) GetListByUserID(ctx context.Context, userID int64) ([]models.Journal, error) {
	var journals []models.Journal

	query := `
		SELECT j.id, j.uid, j.title, j.text, j.created_at, j.updated_at
		FROM journals j
		WHERE user_id = $1
	`

	rows, err := j.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var j models.Journal
		err := rows.Scan(&j.ID, &j.Uid, &j.Title, &j.Text, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		journals = append(journals, j)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}

func (j *journal) Update(ctx context.Context, journal *models.Journal) error {
	query := `
		UPDATE journals
		SET title = $1,
			text = $2,
			mood_id = $3,
			updated_at = $4
		WHERE uid = $5
	`

	result, err := j.pool.Exec(ctx, query, journal.Title, journal.Text, journal.MoodID, time.Now().UTC(), journal.Uid)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}
