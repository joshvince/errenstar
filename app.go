package main

import (
	"context"
	"errenstar/internal/llm"
)

// App struct
type App struct {
	ctx        context.Context
	llmService llm.LLMService
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
}

// CallLLM calls the local LLM model with the given input
func (a *App) CallLLM(input string) string {
	return a.llmService.Ask(a.ctx, input)
}

func (a *App) CancelLLMRequest() string {
	return a.llmService.CancelRequest()
}
