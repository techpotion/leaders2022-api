package service

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileService interface {
	WriteStringToFile(data []byte, filename string) error
}

type fileService struct {
	storagePath string
}

func NewLocalFileService(storagePath string) *fileService {
	return &fileService{storagePath: storagePath}
}

func (s *fileService) WriteStringToFile(data []byte, filename string) error {
	if data == nil {
		return fmt.Errorf("empty data received")
	}

	path := s.storagePath + "/" + filename

	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write data to file: %w", err)
	}

	return nil
}
