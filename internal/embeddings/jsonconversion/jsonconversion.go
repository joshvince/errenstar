package jsonconversion

import (
	"bytes"
	"encoding/json"
	"html"

	"github.com/microcosm-cc/bluemonday"
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
	// Use bluemonday to strip all HTML tags but preserve text content
	p := bluemonday.StrictPolicy()
	cleanText := p.Sanitize(rawString)

	// Decode HTML entities like &#39; to actual characters
	decodedText := html.UnescapeString(cleanText)

	return decodedText
}
