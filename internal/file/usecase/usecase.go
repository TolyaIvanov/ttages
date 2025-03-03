package usecase

import (
	"context"
	"io"
	"ttages/internal/config"
	"ttages/internal/file/DTO"
	"ttages/internal/file/entity"
	"ttages/internal/file/repository/postgres"
	"ttages/internal/file/repository/redis"
)

type FileUsecaseInterface interface {
	Upload(ctx context.Context, req *DTO.FileUpload) (*entity.File, error)
	Download(ctx context.Context, filename string) (io.Reader, error)
	ListFiles(ctx context.Context) ([]entity.File, error)
}

type FileUsecase struct {
	fileRepo  *postgres.FileRepository
	fileCache *redis.FileCache
	cfg       *config.Config
}

func NewFileUsecase(repo *postgres.FileRepository, cache *redis.FileCache, cfg *config.Config) *FileUsecase {
	return &FileUsecase{
		fileRepo:  repo,
		fileCache: cache,
		cfg:       cfg,
	}
}
