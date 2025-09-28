package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Comment struct {
	ID              int32         `json:"id"`
	NewsItemID      sql.NullInt32 `json:"news_item_id"`
	UserID          uuid.NullUUID `json:"user_id"`
	ParentCommentID sql.NullInt32 `json:"parent_comment_id"`
	Content         string        `json:"content"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}
