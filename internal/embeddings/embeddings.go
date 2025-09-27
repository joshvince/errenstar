package embeddings

import (
	"context"
	"errenstar/internal/embeddings/db"
	"errenstar/internal/embeddings/fileops"

	"github.com/philippgille/chromem-go"
)

const (
	rawContentDirectory = "raw_content"
)

type EmbeddingsService struct {
	db *db.EmbeddingsDB
}

func NewEmbeddingsService() *EmbeddingsService {
	return &EmbeddingsService{
		db: db.InitializeDB(),
	}
}

func (service *EmbeddingsService) GetDocumentCount(appContext context.Context) int {
	return service.db.Collection.Count()
}

func (service *EmbeddingsService) FetchContexts(appContext context.Context, question string) []string {
	var response []string

	results, err := service.db.QueryDB(appContext, question)
	if err != nil {
		panic(err)
	}

	handlers := uniqueFileHandlersForResults(results)
	var fileContents []byte

	for _, handler := range handlers {
		fileBytes, err := handler.Read()
		if err == nil {
			fileContents = fileBytes
			response = append(response, string(fileContents))
		}
	}

	return response
}

func uniqueFileHandlersForResults(results []chromem.Result) []fileops.FileHandler {
	var handlers []fileops.FileHandler

	// This creates a Set, effectively. It's memory efficient because the struct{} is 0 bytess
	uniquePaths := make(map[string]struct{})

	for _, result := range results {
		uniquePaths[result.Metadata["original_document"]] = struct{}{}
	}

	for path := range uniquePaths {
		// Ignore if the file cannot be found
		handler, _ := fileops.NewFileHandler(path)

		handlers = append(handlers, *handler)
	}

	return handlers
}
