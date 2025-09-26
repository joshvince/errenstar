package embeddings

import (
	"context"
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
			log.Printf("Processing file: %s", path)

			handler, err := fileops.NewFileHandler(path)
			if err != nil {
				log.Printf("Error creating file handler for %s: %v", path, err)
				return err
			}

			doc := []chromem.Document{generateDocumentFromFile(handler)}
			log.Printf("Generated document with ID: %s", doc[0].ID)

			err = db.collection.AddDocuments(appContext, doc, runtime.NumCPU())
			if err != nil {
				log.Printf("Error adding document %s: %v", doc[0].ID, err)
				return err
			}

			log.Printf("Successfully added document: %s", doc[0].ID)
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

	// Use the full file path as the unique ID
	id := path
	categoryName := filepath.Base(parentDir)

	return id, categoryName
}
