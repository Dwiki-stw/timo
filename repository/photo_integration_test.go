package repository

import (
	"context"
	"database/sql"
	"testing"
	"timo/models"

	"github.com/stretchr/testify/assert"
)

func TestPhotoRepository_Create(t *testing.T) {
	ctx := context.Background()
	repo := NewPhoto(testDB)

	photo := &models.Photo{
		JournalID: 17,
		Url:       "url_test",
	}

	err := repo.Create(ctx, photo)

	assert.NoError(t, err)
	assert.NotZero(t, photo.ID)

	_, _ = testDB.Exec(ctx, `DELETE FROM photos WHERE id = $1`, photo.ID)
}

func TestPhotoRepository_GetByJounalID(t *testing.T) {
	ctx := context.Background()
	repo := NewPhoto(testDB)

	listUrl := []string{"url_1", "url_2", "url_3"}

	for _, url := range listUrl {
		p := &models.Photo{JournalID: 17, Url: url}
		err := repo.Create(ctx, p)
		assert.NoError(t, err)
	}

	photos, err := repo.GetByJournalID(ctx, 1)

	assert.NoError(t, err)

	for i, p := range photos {
		assert.Equal(t, listUrl[i], p.Url)
	}

	defer testDB.Exec(ctx, `DELETE FROM photos`)
}

func TestPhotoRepository_Delete(t *testing.T) {
	tests := []struct {
		name  string
		photo *models.Photo
		isErr bool
	}{
		{
			name:  "no rows",
			isErr: true,
		},
		{
			name:  "success",
			photo: &models.Photo{JournalID: 17, Url: "url_test"},
			isErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := NewPhoto(testDB)

			if tt.isErr {
				err := repo.Delete(ctx, 1)

				assert.Error(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
			} else {
				err := repo.Create(ctx, tt.photo)
				assert.NoError(t, err)

				err = repo.Delete(ctx, tt.photo.ID)
				assert.NoError(t, err)

			}
		})
	}
}
