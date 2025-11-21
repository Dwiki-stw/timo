package models

import "time"

type Photo struct {
	ID        int64     `db:"id"`
	JournalID int64     `db:"journal_id"`
	Url       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}
