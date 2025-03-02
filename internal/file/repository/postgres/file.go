package postgres

import (
	"context"
	"database/sql"

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
		INSERT INTO files (id, name, created_at, updated_at, size)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		file.ID, file.Name, file.CreatedAt, file.UpdatedAt, file.Size)
	return err
}

func (r *FileRepository) GetAll(ctx context.Context) ([]entity.File, error) {
	var files []entity.File
	query := `SELECT * FROM files`
	err := r.db.SelectContext(ctx, &files, query)
	return files, err
}

func (r *FileRepository) GetByName(ctx context.Context, name string) (*entity.File, error) {
	var file entity.File
	query := `SELECT * FROM files WHERE name = $1`
	err := r.db.GetContext(ctx, &file, query, name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &file, err
}
