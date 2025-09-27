package main

import (
	"context"
	"errenstar/internal/embeddings"
	"errenstar/internal/llm"
	"log"
)

// App struct
type App struct {
	ctx               context.Context
	llmService        llm.LLMService
	embeddingsService embeddings.EmbeddingsService
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
	a.embeddingsService = *embeddings.NewEmbeddingsService()

	// Log document count on app startup
	count := a.embeddingsService.GetDocumentCount(ctx)

	log.Printf("App startup - documents in collection: %d", count)
}

// CallLLM calls the local LLM model with the given input
func (a *App) CallLLM(input string) string {
	contexts := a.embeddingsService.FetchContexts(a.ctx, input)

	return a.llmService.Ask(a.ctx, contexts, input)
}

func (a *App) CancelLLMRequest() string {
	return a.llmService.CancelRequest()
}
