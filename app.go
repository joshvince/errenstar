package main

import (
	"context"
	db "errenstar/internal/embeddings/db"
	"errenstar/internal/llm"
	"log"
)

// App struct
type App struct {
	ctx          context.Context
	llmService   llm.LLMService
	embeddingsDB db.EmbeddingsDB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.llmService = *llm.NewLLMService()
	a.embeddingsDB = *db.InitializeDB()

	// Log document count on app startup
	count := a.embeddingsDB.GetDocumentCount(ctx)

	log.Printf("App startup - documents in collection: %d", count)
}

// CallLLM calls the local LLM model with the given input
func (a *App) CallLLM(input string) string {
	contexts := a.embeddingsDB.QueryDB(a.ctx, input)

	return a.llmService.Ask(a.ctx, contexts, input)
}

func (a *App) CancelLLMRequest() string {
	return a.llmService.CancelRequest()
}
