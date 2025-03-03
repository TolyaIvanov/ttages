package gRPC

import (
	"context"
	"io"
	"time"
	"ttages/internal/file/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ttages/internal/file/DTO"
	"ttages/internal/file/usecase"
	"ttages/proto/pb"
)

type FileHandler struct {
	pb.UnimplementedFileServiceServer
	uc          *usecase.FileUsecase
	uploadSem   chan struct{}
	downloadSem chan struct{}
	listSem     chan struct{}
}

func NewFileHandler(uc *usecase.FileUsecase) *FileHandler {
	return &FileHandler{
		uc:          uc,
		uploadSem:   make(chan struct{}, 10),
		downloadSem: make(chan struct{}, 10),
		listSem:     make(chan struct{}, 100),
	}
}

func (h *FileHandler) Upload(stream pb.FileService_UploadServer) error {
	h.uploadSem <- struct{}{}
	defer func() { <-h.uploadSem }()

	var fileData = &DTO.FileUpload{}

	// Читаем первый чанк (содержит имя файла)
	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	fileData.Name = req.GetFilename()

	// Цикл приема чанков
	for {
		fileData.Chunk = append(fileData.Chunk, req.GetChunk()...)

		req, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	file, err := h.uc.Upload(stream.Context(), fileData)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&pb.UploadResponse{
		Id:   file.ID,
		Size: uint32(file.Size),
	})
}

func (h *FileHandler) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	h.downloadSem <- struct{}{}
	defer func() { <-h.downloadSem }()

	reader, err := h.uc.Download(stream.Context(), req.GetFilename())
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&pb.DownloadResponse{Chunk: buf[:n]}); err != nil {
			return err
		}
	}
	return nil
}

func (h *FileHandler) ListFiles(ctx context.Context, _ *pb.ListRequest) (*pb.ListResponse, error) {
	h.listSem <- struct{}{}
	defer func() { <-h.listSem }()

	files, err := h.uc.ListFiles(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if files == nil {
		files = []entity.File{}
	}

	pbFiles := make([]*pb.FileInfo, 0, len(files))
	for _, f := range files {
		pbFiles = append(pbFiles, &pb.FileInfo{
			Filename:  f.Name,
			CreatedAt: f.CreatedAt.Format(time.RFC3339),
			UpdatedAt: f.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &pb.ListResponse{Files: pbFiles}, nil
}
