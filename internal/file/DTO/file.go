package DTO

import "time"

type FileInfo struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileUpload struct {
	Chunk []byte
	Name  string
}

type FileDownload struct {
	Chunk []byte
}
