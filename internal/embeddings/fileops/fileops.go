package fileops

import (
	"io"
	"os"
)

// FileHandler represents a file manipulation handler
type FileHandler struct {
	file *os.File
}

// NewFileHandler creates a new file handler and opens the file
func NewFileHandler(filePath string) (*FileHandler, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	return &FileHandler{file: file}, nil
}

// Close the file
func (fh *FileHandler) Close() error {
	err := fh.file.Close()

	if err != nil {
		return err
	}

	return nil
}

// Reads the entire file contents into memory
func (fh *FileHandler) Read() ([]byte, error) {
	contents, err := io.ReadAll(fh.file)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
