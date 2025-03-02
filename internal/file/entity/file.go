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
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Size      int64     `db:"size"`
}
