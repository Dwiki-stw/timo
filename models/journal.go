package models

import "time"

type Journal struct {
	ID        int64     `db:"id"`
	Uid       string    `db:"uid"`
	UserID    int64     `db:"user_id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	MoodID    int64     `db:"mood_id"`
	MoodLabel string    `db:"mood_label"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
