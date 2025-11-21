package repository

import (
	"context"
	"database/sql"
	"testing"
	"timo/models"

	"github.com/stretchr/testify/assert"
)

func TestJournalRepository_Create(t *testing.T) {
	ctx := context.Background()
	journal := &models.Journal{
		UserID: 14,
		Title:  "title test",
		Text:   "text test",
		MoodID: 1,
	}

	repo := NewJournal(testDB)
	err := repo.Create(ctx, journal)

	assert.NoError(t, err)
	assert.NotZero(t, journal.ID)
	assert.NotEmpty(t, journal.Uid)

	_, _ = testDB.Exec(ctx, `DELETE FROM journals WHERE id = $1`, journal.ID)
}

func TestJournalRepository_GetByID(t *testing.T) {
	ctx := context.Background()
	repo := NewJournal(testDB)

	expectedJournal := &models.Journal{
		UserID: 14,
		Title:  "title test",
		Text:   "text test",
		MoodID: 1,
	}

	err := repo.Create(ctx, expectedJournal)
	assert.NoError(t, err)

	journal, err := repo.GetByID(ctx, expectedJournal.Uid)

	assert.NoError(t, err)
	assert.Equal(t, expectedJournal.Uid, journal.Uid)
	assert.Equal(t, expectedJournal.Text, journal.Text)

	_, _ = testDB.Exec(ctx, `DELETE FROM journals WHERE id = $1`, journal.ID)
}

func TestJournalRepository_Update(t *testing.T) {
	ctx := context.Background()
	repo := NewJournal(testDB)

	oldJournal := &models.Journal{
		UserID: 14,
		Title:  "title test",
		Text:   "text test",
		MoodID: 1,
	}

	err := repo.Create(ctx, oldJournal)
	assert.NoError(t, err)

	newJournal := &models.Journal{
		ID:     oldJournal.ID,
		Uid:    oldJournal.Uid,
		UserID: 14,
		Title:  "title update",
		Text:   "text update",
		MoodID: 1,
	}

	err = repo.Update(ctx, newJournal)
	assert.NoError(t, err)

	journal, err := repo.GetByID(ctx, oldJournal.Uid)

	assert.NoError(t, err)
	assert.Equal(t, journal.Text, newJournal.Text)
	assert.NotEqual(t, journal.Text, oldJournal.Text)

	_, _ = testDB.Exec(ctx, `DELETE FROM journals WHERE id = $1`, journal.ID)
}

func TestJournalRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		journal *models.Journal
		err     bool
	}{
		{
			name: "no rows",
			err:  true,
		},
		{
			name:    "success",
			journal: &models.Journal{UserID: 14, Title: "title test", Text: "text test", MoodID: 1},
			err:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := NewJournal(testDB)

			if tt.err {
				err := repo.Delete(ctx, "550e8400-e29b-41d4-a716-446655440000")

				assert.Error(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
			} else {
				err := repo.Create(ctx, tt.journal)
				assert.NoError(t, err)

				err = repo.Delete(ctx, tt.journal.Uid)
				assert.NoError(t, err)

				journal, err := repo.GetByID(ctx, tt.journal.Uid)
				assert.Error(t, err)
				assert.Equal(t, sql.ErrNoRows, err)
				assert.Nil(t, journal)
			}
		})
	}
}

func TestJournalRepository_GetListByUserID(t *testing.T) {
	ctx := context.Background()
	repo := NewJournal(testDB)

	listText := []string{"text 1", "text 2", "text 3"}

	for _, text := range listText {
		j := &models.Journal{
			UserID: 14,
			Title:  "title test",
			Text:   text,
			MoodID: 1,
		}
		err := repo.Create(ctx, j)
		assert.NoError(t, err)
	}

	journals, err := repo.GetListByUserID(ctx, 14)

	assert.NoError(t, err)
	for i, j := range journals {
		assert.Equal(t, listText[i], j.Text)
	}
}
