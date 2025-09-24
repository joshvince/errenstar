package llm

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type LLMService struct {
	cancelFunc context.CancelFunc
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

func (service *LLMService) Ask(appContext context.Context, input string) string {
	printInput(input)
	ctx, cancel := context.WithCancel(appContext)
	service.cancelFunc = cancel

	return callCLIModel(ctx, input)
}

func (service *LLMService) CancelRequest() string {
	if service.cancelFunc != nil {
		service.cancelFunc()
		service.cancelFunc = nil
	}

	return "Cancelled the request"
}

func printInput(input string) {
	fmt.Print("Receiving input: ", input)
}

func callCLIModel(ctx context.Context, input string) string {
	cmd := exec.CommandContext(ctx, "llm", "-m", "mlx-community/Llama-3.2-3B-Instruct-4bit", input)

	output, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf("Error calling LLM: %v", err)
	}

	return strings.TrimSpace(string(output))
}
