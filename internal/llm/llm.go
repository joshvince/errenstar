package llm

import (
	"context"
	"fmt"
	"html/template"
	"os/exec"
	"strings"
)

type LLMService struct {
	cancelFunc context.CancelFunc
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

func (service *LLMService) CancelRequest() string {
	if service.cancelFunc != nil {
		service.cancelFunc()
		service.cancelFunc = nil
	}

	return "Cancelled the request"
}

func (service *LLMService) Ask(appContext context.Context, contexts []string, question string) string {
	fmt.Print("Received input: ", question)

	ctx, cancel := context.WithCancel(appContext)
	service.cancelFunc = cancel

	return callCLIModel(ctx, contexts, question)
	// You can replace this with this line if you want to call a remote server
	// return askRemoteOllamaModel(ctx, contexts, question)
}

func callCLIModel(ctx context.Context, userContexts []string, input string) string {
	systemPrompt := cliSystemPromptWithContext(userContexts)
	userInput := input

	fmt.Printf("Gave the model this system prompt: %s", systemPrompt)
	fmt.Printf("Asked the model: %s", userInput)

	cmd := exec.CommandContext(
		ctx, "llm", userInput,
		"-m", "mlx-community/Llama-3.2-3B-Instruct-4bit",
		"-o", "temperature", "0.1", "-s", systemPrompt,
	)

	output, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf("Error calling LLM: %v", err)
	}

	return strings.TrimSpace(string(output))
}

func cliSystemPromptWithContext(contexts []string) string {
	systemPromptTemplate := template.Must(template.New("system_prompt_with_context").Parse(`
You are a knowledge assistant for the fictional world of Sonovem, a Dungeons & Dragons campaign by Lloyd Morgan.

The player characters in this campaign will be referred to by the user as either the party, the team, the Royaum Rippers, or even just the Rippers.
If you receive a question that asks about "us" or other similar informal term, you can infer that this means the party or player characters.

CRITICAL RULES - YOU MUST FOLLOW THESE EXACTLY:
1. ONLY answer based on information explicitly provided in the context below
2. NEVER add details, speculation, or creative elements not in the context
3. NEVER make assumptions or fill in gaps with your own knowledge
4. If information is not in the context, respond with "I don't have information about that in the knowledge base"
5. Use only facts directly stated in the provided context
6. Keep responses concise and factual

Answer in a neutral, journalistic tone. Do not embellish, interpret, or expand beyond what is explicitly written in the context.
Answer ONLY using the information provided in the context below. If the context doesn't contain relevant information for the question, respond with "I don't have information about that in the knowledge base."

The following context contains information from the Sonovem knowledge base:

<context>
{{range $i, $context := .}}{{if $i}}	----
{{end}}	{{$context}}
{{end}}</context>`))

	sb := &strings.Builder{}
	err := systemPromptTemplate.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}
	return sb.String()
}
