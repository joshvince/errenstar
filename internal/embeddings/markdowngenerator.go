package embeddings

import (
	"errenstar/internal/embeddings/fileops"
	jsonconversion "errenstar/internal/embeddings/json_conversion"
	"strings"
)

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
