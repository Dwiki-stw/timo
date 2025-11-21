package domain

import (
	"context"
	"timo/models"
)

type PhotoRepository interface {
	GetByJournalID(ctx context.Context, journalID int64) ([]models.Photo, error)
	Create(ctx context.Context, photo *models.Photo) error
	Delete(ctx context.Context, id int64) error
}
