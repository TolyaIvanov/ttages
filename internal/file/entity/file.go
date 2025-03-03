package entity

import (
	"errors"
	"time"
)

var (
	ErrFileExists   = errors.New("file already exists")
	ErrFileNotFound = errors.New("file not found")
)

type File struct {
	ID        string    `db:"id"`
	Name      string    `db:"filename"`
	Path      string    `db:"path"`
	Size      int64     `db:"size"`
	MimeType  string    `db:"mime_type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
