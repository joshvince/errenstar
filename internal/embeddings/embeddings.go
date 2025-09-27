package embeddings

import (
	"context"
	"errenstar/internal/embeddings/db"
	"log"
	"strings"
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
	results, err := service.db.QueryDB(appContext, question)
	if err != nil {
		panic(err)
	}

	var response []string

	// TODO: we need to fetch the unique original_document values from the metadata
	// Then use that to fetch some files using fileops
	// And _that_ is the slice we should send as context

	for i, res := range results {
		content := strings.TrimPrefix(res.Content, "search_document: ")

		log.Printf("Document %d (similarity: %f): \"%s\"\n", i+1, res.Similarity, content)

		response = append(response, content)
	}

	return response
}
