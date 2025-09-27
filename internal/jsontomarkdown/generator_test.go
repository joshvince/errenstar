package jsontomarkdown

import (
	"errenstar/internal/embeddings/fileops"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempDirAndFileHandler(t *testing.T) *fileops.FileHandler {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "character_example.json")

	inputData, err := os.ReadFile("testdata/characters/character_example.json")
	require.NoError(t, err)
	err = os.WriteFile(inputPath, inputData, 0644)
	require.NoError(t, err)

	testFileHandler, err := fileops.NewFileHandler(inputPath)
	require.NoError(t, err)

	return testFileHandler
}

func TestCreateMarkdownFromJSONFile(t *testing.T) {
	inputFileHandler := createTempDirAndFileHandler(t)

	outputFileHandler, err := createMarkdownFromJSONFile(inputFileHandler)
	require.NoError(t, err)

	expectedOutputHandler, err := fileops.NewFileHandler("testdata/characters/expected_output.md")
	require.NoError(t, err)

	expectedOutput, _ := expectedOutputHandler.Read()
	actualOutput, err := outputFileHandler.Read()
	require.NoError(t, err)

	assert.Equal(t, string(actualOutput), string(expectedOutput))
}
