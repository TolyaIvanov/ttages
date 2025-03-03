package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"ttages/internal/file/entity"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *entity.File) error {
	query := `
		INSERT INTO files (id, filename, path, size, mime_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		file.ID, file.Name, file.Path, file.Size, file.MimeType, file.CreatedAt, file.UpdatedAt)

	return err
}

func (r *FileRepository) GetAll(ctx context.Context) ([]entity.File, error) {
	var files []entity.File
	query := `SELECT id, filename, size, mime_type, created_at, updated_at FROM files`
	err := r.db.SelectContext(ctx, &files, query)

	return files, err
}

func (r *FileRepository) GetByName(ctx context.Context, name string) (*entity.File, error) {
	var file entity.File
	query := `
        SELECT id, filename, path, size, created_at, updated_at 
        FROM files 
        WHERE filename = $1`
	err := r.db.GetContext(ctx, &file, query, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &file, err
}
