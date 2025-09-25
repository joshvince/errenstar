package embeddings

import (
	"errenstar/internal/embeddings/fileops"
	"errenstar/internal/embeddings/jsonconversion"
	"io/fs"
	"path/filepath"
	"strings"
)

func ConvertDirectoriesToMarkdown(path string) error {
	err := filepath.WalkDir(path, visit)
	if err != nil {
		return err
	}

	return nil
}

func visit(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !entry.IsDir() && strings.HasSuffix(path, ".json") {
		handler, err := fileops.NewFileHandler(path)
		if err != nil {
			return err
		}

		_, err = createMarkdownFromJSONFile(handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func createMarkdownFromJSONFile(handler *fileops.FileHandler) (*fileops.FileHandler, error) {
	var contents []byte

	contents, err := handler.Read()
	if err != nil {
		return nil, err
	}

	contents, err = jsonconversion.ExtractMarkdownFromJSON(contents)
	if err != nil {
		return nil, err
	}

	fileName := strings.TrimSuffix(handler.GetPath(), ".json") + ".md"

	outputHandler, err := fileops.Write(fileName, contents)
	if err != nil {
		return nil, err
	}

	return outputHandler, nil
}
