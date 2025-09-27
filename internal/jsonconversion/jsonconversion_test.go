package jsonconversion

import (
	"errenstar/internal/embeddings/fileops"
	"regexp"
	"strings"
	"testing"
)

func readTestFileContents(t *testing.T) []byte {
	fixturePath := "../embeddings/testdata/characters/character_example.json"
	handler, err := fileops.NewFileHandler(fixturePath)

	if err != nil {
		t.Fatalf("Failed to open file fixture: %s: %v", fixturePath, err)
	}

	t.Cleanup(func() {
		defer handler.Close()
	})

	contents, err := handler.Read()
	if err != nil {
		t.Fatalf("Failed to read file fixture: %s: %v", fixturePath, err)
	}

	return contents
}

func TestConvertKankaToStruct(t *testing.T) {
	testFileContents := readTestFileContents(t)

	_, err := convertKankaToStruct(testFileContents)
	if err != nil {
		t.Errorf("Failed to convert to string: %v", err)
	}
}

func TestExtractMarkdownFromJSON_NoHistory(t *testing.T) {
	testFileContents := readTestFileContents(t)

	result, err := ExtractMarkdownFromJSON(testFileContents)
	if err != nil {
		t.Errorf("Failed to write to string %v", err)
	}

	if strings.Contains(strings.ToLower(string(result)), "history") {
		t.Errorf("Result should omit anything inside HTML tag attributes, but got: %s", result)
	}
}

func TestExtractMarkdownFromJSON_NoHTMLTags(t *testing.T) {
	testFileContents := readTestFileContents(t)

	result, err := ExtractMarkdownFromJSON(testFileContents)
	if err != nil {
		t.Errorf("Failed to write to string %v", err)
	}

	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	if htmlTagRegex.MatchString(string(result)) {
		t.Errorf("Result should not contain HTML tags, but got: %s", result)
	}
}
