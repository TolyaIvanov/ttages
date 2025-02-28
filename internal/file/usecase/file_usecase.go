package usecase

import (
	"ttages/internal/file/DTO"
	"ttages/internal/file/entity"
	"ttages/internal/file/repository"
)

type FileUsecase struct {
	repo repository.FileRepository
}

func NewFileUsecase(repo repository.FileRepository) *FileUsecase {
	return &FileUsecase{repo: repo}
}

func (uc *FileUsecase) UploadFile(fileDTO DTO.FileDTO) error {
}

func (uc *FileUsecase) GetFiles() ([]DTO.FileDTO, error) {
}
