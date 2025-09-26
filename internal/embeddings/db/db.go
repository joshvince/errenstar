package embeddings

import (
	"context"
	"errenstar/internal/embeddings/fileops"
	"path/filepath"
	"runtime"

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
	handler, err := fileops.NewFileHandler(rawContentDirectory + "/characters/crispin-tendies.md")
	if err != nil {
		panic(err)
	}

	docs := []chromem.Document{generateDocumentFromFile(handler)}

	err = db.collection.AddDocuments(appContext, docs, runtime.NumCPU())
	if err != nil {
		panic(err)
	}

}

func generateDocumentFromFile(handler *fileops.FileHandler) chromem.Document {
	markdown, err := handler.Read()
	if err != nil {
		panic(err)
	}

	content := "search_document" + string(markdown)

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
