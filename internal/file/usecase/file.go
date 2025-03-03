package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	_ "ttages/internal/config"

	"ttages/internal/file/DTO"
	"ttages/internal/file/entity"
	_ "ttages/internal/file/repository/postgres"
	_ "ttages/internal/file/repository/redis"
)

func (uc *FileUsecase) Upload(ctx context.Context, req *DTO.FileUpload) (*entity.File, error) {
	log.Println("UC Upload")

	existing, err := uc.fileRepo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, entity.ErrFileExists
	}
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(uc.cfg.StoragePath, req.Name)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.Write(req.Chunk); err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}

	now := time.Now()
	newFile := &entity.File{
		ID:        id.String(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
		Size:      int64(len(req.Chunk)),
	}

	if err := uc.fileRepo.Create(ctx, newFile); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return newFile, nil
}

func (uc *FileUsecase) Download(ctx context.Context, filename string) (io.Reader, error) {
	log.Println("UC Download")

	file, err := os.Open(filepath.Join(uc.cfg.StoragePath, filename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, entity.ErrFileNotFound
		}
		return nil, err
	}
	return file, nil
}

func (uc *FileUsecase) ListFiles(ctx context.Context) ([]entity.File, error) {
	log.Println("UC list")

	cached, err := uc.fileCache.GetFiles(ctx)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	files, err := uc.fileRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if err := uc.fileCache.SetFiles(ctx, files); err != nil {
		return files, nil
	}
	return files, nil
}
