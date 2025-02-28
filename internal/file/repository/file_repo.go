package repository

import "ttages/internal/file/entity"

type FileRepository interface {
	Save(file entity.File) error
	FindAll() ([]entity.File, error)
	FindByID(id string) (*entity.File, error)
}
