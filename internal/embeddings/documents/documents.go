package embeddings

import (
	"errenstar/internal/embeddings/fileops"
	"fmt"
	"log"
	"path/filepath"

	"github.com/philippgille/chromem-go"
	"github.com/tmc/langchaingo/textsplitter"
)

func GenerateDocumentsFromFile(handler *fileops.FileHandler) []chromem.Document {
	fileContents, err := handler.Read()
	if err != nil {
		panic(err)
	}

	var documents []chromem.Document

	filePath, category := getInfoFromFilePath(handler)

	chunks := splitMarkdownIntoChunks(string(fileContents))

	for i, chunk := range chunks {
		content := "search_document: " + chunk

		doc := chromem.Document{
			ID:       filePath + "_" + fmt.Sprint(i),
			Metadata: map[string]string{"category": category, "original_document": filePath},
			Content:  content,
		}

		log.Printf("Generated document with ID: %s", doc.ID)

		documents = append(documents, doc)
	}

	return documents
}

func splitMarkdownIntoChunks(fileContents string) []string {
	splitter := textsplitter.NewRecursiveCharacter()

	chunks, err := splitter.SplitText(fileContents)
	if err != nil {
		panic(err)
	}

	return chunks
}

func getInfoFromFilePath(handler *fileops.FileHandler) (string, string) {
	path := handler.GetPath()
	parentDir := filepath.Dir(path)

	categoryName := filepath.Base(parentDir)

	// Use the full file path as the unique ID
	return path, categoryName
}
