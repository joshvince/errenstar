package fileops

import (
	"os"
	"testing"
)

func createTestFileHandler(t *testing.T, fileContent string) *FileHandler {
	// Create a temporary file for testing and open it
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}

	_, err = tmpFile.WriteString(fileContent)
	if err != nil {
		t.Fatal("Failed to write test content:", err)
	}

	t.Cleanup(func() {
		defer os.Remove(tmpFile.Name()) // Remove the file
	})

	tmpFile.Close()

	handler, err := NewFileHandler(tmpFile.Name())
	if err != nil {
		t.Errorf("Failed to open existing file: %v", err)
	}

	t.Cleanup(func() {
		handler.Close()
	})

	return handler
}

func TestFileHandler_Open_Close(t *testing.T) {
	handler := createTestFileHandler(t, "Hello")

	// Test closing the file
	err := handler.Close()
	if err != nil {
		t.Errorf("Failed to close file: %v", err)
	}
}

func TestFileHandler_Read(t *testing.T) {
	fileContent := "Hello, World!"
	handler := createTestFileHandler(t, fileContent)

	// Test reading the file
	data, err := handler.Read()
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if string(data) != fileContent {
		t.Errorf("Expected content %q, got %q", fileContent, string(data))
	}
}
