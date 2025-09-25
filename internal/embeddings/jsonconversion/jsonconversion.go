package jsonconversion

import (
	"bytes"
	"encoding/json"
	"regexp"
)

type KankaJSON struct {
	Name  string `json:"name"`
	Entry string `json:"entry"`
}

func ExtractMarkdownFromJSON(fileInput []byte) ([]byte, error) {
	var output []byte

	kankaStruct, err := convertKankaToStruct(fileInput)
	if err != nil {
		return output, err
	}

	var buffer bytes.Buffer

	buffer.WriteString("# Name:\n")
	buffer.WriteString(kankaStruct.Name)
	buffer.WriteString("\n\n")
	buffer.WriteString("# Details:\n")
	buffer.WriteString(stripHTMLTags(kankaStruct.Entry))
	output = buffer.Bytes()

	return output, nil
}

func convertKankaToStruct(fileInput []byte) (KankaJSON, error) {
	kankaObj := KankaJSON{}

	result, err := fileInputToStruct(fileInput, kankaObj)
	if err != nil {
		return kankaObj, err
	}

	return result, nil
}

func fileInputToStruct(fileInput []byte, kankaObj KankaJSON) (KankaJSON, error) {
	if err := json.Unmarshal(fileInput, &kankaObj); err != nil {
		return kankaObj, err
	}

	return kankaObj, nil
}

func stripHTMLTags(rawString string) string {
	// Strips the contents of the tags but not the text inside
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	cleanText := htmlTagRegex.ReplaceAllString(rawString, "")

	return cleanText
}
