package embeddings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func readTestData(t *testing.T) string {
	inputData, err := os.ReadFile("../testdata/locations/errenstar.md")
	require.NoError(t, err)

	return string(inputData)
}

func TestSplitMarkdownIntoChunks(t *testing.T) {
	inputData := readTestData(t)

	chunks := splitMarkdownIntoChunks(inputData)

	if len(chunks) != 5 {
		t.Errorf("SplitMarkdownIntoChunks split into %v instead of 3 chunks", len(chunks))
	}
}
