package repository

import (
	"github.com/jmoiron/sqlx"
	"ttages/internal/file/entity"
)

type PostgresRepository struct {
	DB *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

func (r *PostgresRepository) Save(file entity.File) error {
}

func (r *PostgresRepository) FindAll() ([]entity.File, error) {
}

func (r *PostgresRepository) FindByID(id string) (*entity.File, error) {
}
