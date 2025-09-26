package embeddings

import (
	"context"
	"errenstar/internal/embeddings/fileops"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/philippgille/chromem-go"
)

const (
	embeddingModel = "nomic-embed-text"
	embeddingURL   = "http://192.168.1.160:11434/api"
	dbPath         = "./db"
)

type EmbeddingsDB struct {
	db         *chromem.DB
	collection *chromem.Collection
}

func InitializeDB() *EmbeddingsDB {
	var db *chromem.DB
	var err error

	// Always use NewPersistentDB - it will create if doesn't exist, or open if it does
	db, err = chromem.NewPersistentDB(dbPath, false)
	if err != nil {
		panic(err)
	}

	collection, err := db.GetOrCreateCollection("Sonovem", nil, chromem.NewEmbeddingFuncOllama(embeddingModel, embeddingURL))
	if err != nil {
		panic(err)
	}

	return &EmbeddingsDB{
		db:         db,
		collection: collection,
	}
}

func (db *EmbeddingsDB) SeedDB(appContext context.Context) {
	rawContentDirectory := "raw_content"

	docs := loadAllMarkdown(rawContentDirectory)

	err := db.collection.AddDocuments(appContext, docs, runtime.NumCPU())
	if err != nil {
		panic(err)
	}

}

func (db *EmbeddingsDB) QueryDB(appContext context.Context, question string) []string {
	query := "search_query: " + question
	var response []string

	docRes, err := db.collection.Query(appContext, query, 2, nil, nil)
	if err != nil {
		panic(err)
	}

	for i, res := range docRes {
		content := string(res.Content)

		if embeddingModel == "nomic-embed-text" {
			// This prefix is specific to the "nomic-embed-text" model.
			content = strings.TrimPrefix(res.Content, "search_document: ")
		}
		log.Printf("Document %d (similarity: %f): \"%s\"\n", i+1, res.Similarity, content)

		response = append(response, content)
	}

	return response
}

func loadAllMarkdown(directory string) []chromem.Document {
	var docs []chromem.Document
	// TODO: we should eventually load all the docs here
	firstFilePath := directory + "/characters/crispin-tendies.md"
	secondFilePath := directory + "/locations/acquenti.md"

	paths := []string{firstFilePath, secondFilePath}
	for _, path := range paths {
		handler, err := fileops.NewFileHandler(path)
		if err != nil {
			panic(err)
		}

		doc := generateDocumentFromFile(handler)
		docs = append(docs, doc)
	}
	return docs
}

func generateDocumentFromFile(handler *fileops.FileHandler) chromem.Document {
	markdown, err := handler.Read()
	if err != nil {
		panic(err)
	}

	content := "search_document: " + string(markdown)

	id, category := getInfoFromFilePath(handler)

	return chromem.Document{
		ID:       id,
		Metadata: map[string]string{"category": category},
		Content:  content,
	}
}

func getInfoFromFilePath(handler *fileops.FileHandler) (string, string) {
	path := handler.GetPath()
	parentDir := filepath.Dir(path)

	id := filepath.Dir(path)
	categoryName := filepath.Base(parentDir)

	return id, categoryName
}
