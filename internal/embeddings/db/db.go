package embeddings

import (
	"context"
	documents "errenstar/internal/embeddings/documents"
	"errenstar/internal/embeddings/fileops"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/philippgille/chromem-go"
)

const (
	embeddingModel      = "nomic-embed-text"
	embeddingURL        = "http://192.168.1.160:11434/api"
	dbPath              = "./db"
	rawContentDirectory = "raw_content"
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

func (db *EmbeddingsDB) GetDocumentCount(appContext context.Context) int {
	return db.collection.Count()
}

func (db *EmbeddingsDB) SeedDB(appContext context.Context) {
	// Get initial document count
	initialCount := db.collection.Count()

	log.Printf("Initial document count: %d", initialCount)

	loadMarkdown := func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() && strings.HasSuffix(path, ".md") {
			handler, err := fileops.NewFileHandler(path)
			if err != nil {
				log.Printf("Error creating file handler for %s: %v", path, err)
				return err
			}

			docs := documents.GenerateDocumentsFromFile(handler)

			err = db.collection.AddDocuments(appContext, docs, runtime.NumCPU())
			if err != nil {
				log.Printf("Error adding document %v", err)
				return err
			}
		}
		return nil
	}

	err := filepath.WalkDir(rawContentDirectory, loadMarkdown)
	if err != nil {
		panic(err)
	}

	// Get final document count
	finalCount := db.collection.Count()
	log.Printf("Final document count: %d (added %d documents)", finalCount, finalCount-initialCount)
}

func (db *EmbeddingsDB) QueryDB(appContext context.Context, question string) []string {
	query := "search_query: " + question
	var response []string

	docRes, err := db.collection.Query(appContext, query, 10, nil, nil)
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
