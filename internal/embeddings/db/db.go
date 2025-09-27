package db

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
	embeddingURL        = "http://127.0.0.1:11434/api"
	dbPath              = "./db"
	rawContentDirectory = "raw_content"
)

type EmbeddingsDB struct {
	Db         *chromem.DB
	Collection *chromem.Collection
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
		Db:         db,
		Collection: collection,
	}
}

func (db *EmbeddingsDB) SeedDB(appContext context.Context) {
	// Get initial document count
	initialCount := db.Collection.Count()

	log.Printf("Initial document count: %d", initialCount)

	loadMarkdown := func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() && strings.HasSuffix(path, ".md") {
			handler, err := fileops.NewFileHandler(path)
			if err != nil {
				log.Printf("Error creating file handler for %s: %v", path, err)
				return err
			}

			docs := documents.GenerateDocumentsFromFile(handler)

			err = db.Collection.AddDocuments(appContext, docs, runtime.NumCPU())
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
	finalCount := db.Collection.Count()
	log.Printf("Final document count: %d (added %d documents)", finalCount, finalCount-initialCount)
}

func (db *EmbeddingsDB) QueryDB(appContext context.Context, question string) ([]chromem.Result, error) {
	query := "search_query: " + question

	docRes, err := db.Collection.Query(appContext, query, 10, nil, nil)
	if err != nil {
		return nil, err
	}

	return docRes, err
}
