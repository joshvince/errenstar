package jsonconversion

import (
	"encoding/json"
	"regexp"
	"strings"
)

type KankaJSON struct {
	Name  string `json:"name"`
	Entry string `json:"entry"`
}

func ExtractMarkdownFromJSON(fileInput []byte) (string, error) {
	var outputString string

	kankaStruct, err := convertKankaToStruct(fileInput)
	if err != nil {
		return outputString, err
	}

	var builder strings.Builder

	builder.WriteString("# Name: \n")
	builder.WriteString(kankaStruct.Name)
	builder.WriteString("\n\n")
	builder.WriteString("# Details: \n")

	builder.WriteString(stripHTMLTags(kankaStruct.Entry))

	return builder.String(), nil
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
