package domain

import (
	"context"
	"timo/models"
)

type JournalRepository interface {
	GetListByUserID(ctx context.Context, userID int64) ([]models.Journal, error)
	GetByID(ctx context.Context, uid string) (*models.Journal, error)
	Create(ctx context.Context, journal *models.Journal) error
	Update(ctx context.Context, journal *models.Journal) error
	Delete(ctx context.Context, uid string) error
}
