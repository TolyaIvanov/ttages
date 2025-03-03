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

	existingFile, err := uc.fileRepo.GetByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingFile != nil {
		return nil, entity.ErrFileExists
	}

	// Сохраняем файл на диск
	filePath := filepath.Join(uc.cfg.StoragePath, req.Name)
	if err := os.WriteFile(filePath, req.Chunk, 0777); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	newFile := &entity.File{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Path:      "/storage/" + req.Name,
		Size:      int64(len(req.Chunk)),
		MimeType:  "application/octet-stream",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.fileRepo.Create(ctx, newFile)
	if err != nil {
		return nil, err
	}

	return newFile, nil
}

func (uc *FileUsecase) Download(ctx context.Context, filename string) (io.Reader, error) {
	filePath := filepath.Join(uc.cfg.StoragePath, filename)
	log.Printf("Attempting to open file: %s", filePath) // Логируем путь

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
