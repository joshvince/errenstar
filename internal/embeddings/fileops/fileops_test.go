package fileops

import (
	"os"
	"testing"
)

func CreateTestFileHandler(t *testing.T, fileContent string) *FileHandler {
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
	handler := CreateTestFileHandler(t, "Hello")

	// Test closing the file
	err := handler.Close()
	if err != nil {
		t.Errorf("Failed to close file: %v", err)
	}
}

func TestFileHandler_Read(t *testing.T) {
	fileContent := "Hello, World!"
	handler := CreateTestFileHandler(t, fileContent)

	// Test reading the file
	data, err := handler.Read()
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}

	if string(data) != fileContent {
		t.Errorf("Expected content %q, got %q", fileContent, string(data))
	}
}

func TestFileHandler_Write(t *testing.T) {
	path := "test_file.md"
	contentString := "Hello, World!"

	// Clean up the test file after the test
	t.Cleanup(func() {
		os.Remove(path)
	})

	handler, err := Write(path, []byte(contentString))
	if err != nil {
		t.Errorf("Failed to write the file: %v", err)
	}

	outputFileContents, err := handler.Read()
	if err != nil {
		t.Errorf("Failed to read the written file: %v", err)
	}

	outputContentsString := string(outputFileContents)
	if outputContentsString != contentString {
		t.Errorf("The written file did not match the input: %s: %v", outputContentsString, err)
	}

	// Close the handler
	handler.Close()
}
