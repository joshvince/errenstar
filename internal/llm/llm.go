package llm

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	ollamaBaseURL = "http://192.168.1.160:11434/v1"
	llmModel      = "llama3:8b"
)

type LLMService struct {
	cancelFunc context.CancelFunc
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

func (service *LLMService) Ask(appContext context.Context, contexts []string, question string) string {
	fmt.Print("Received input: ", question)

	ctx, cancel := context.WithCancel(appContext)
	service.cancelFunc = cancel

	return askRemoteOllamaModel(ctx, contexts, question)
}

func (service *LLMService) CancelRequest() string {
	if service.cancelFunc != nil {
		service.cancelFunc()
		service.cancelFunc = nil
	}

	return "Cancelled the request"
}

var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
You are a knowledge assistant for the fictional world of Sonovem, a Dungeons & Dragons campaign by Lloyd Morgan.

CRITICAL RULES - YOU MUST FOLLOW THESE EXACTLY:
1. ONLY answer based on information explicitly provided in the context below
2. NEVER add details, speculation, or creative elements not in the context
3. NEVER make assumptions or fill in gaps with your own knowledge
4. If information is not in the context, respond with "I don't have information about that in the knowledge base"
5. Use only facts directly stated in the provided context
6. Keep responses concise and factual

Answer in a neutral, journalistic tone. Do not embellish, interpret, or expand beyond what is explicitly written in the context.
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
Answer ONLY using the information provided in the context below. If the context doesn't contain relevant information for the question, respond with "I don't have information about that in the knowledge base."

The following context contains information from the Sonovem knowledge base:

<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}
`))

func askRemoteOllamaModel(ctx context.Context, contexts []string, question string) string {
	// We can use the OpenAI client because Ollama is compatible with OpenAI's API.
	openAIClient := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    ollamaBaseURL,
		HTTPClient: http.DefaultClient,
	})

	sb := &strings.Builder{}
	err := systemPromptTpl.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Gave the model this system prompt and context: %s", sb.String())
	fmt.Printf("Asked the model: %s", question)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: sb.String(),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "Question: " + question,
		},
	}

	res, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       llmModel,
		Messages:    messages,
		Temperature: 0.1, // Low temperature for more deterministic, factual responses
	})
	if err != nil {
		panic(err)
	}

	reply := res.Choices[0].Message.Content
	reply = strings.TrimSpace(reply)

	return reply
}

func callCLIModel(ctx context.Context, input string) string {
	cmd := exec.CommandContext(ctx, "llm", "-m", "mlx-community/Llama-3.2-3B-Instruct-4bit", input)

	output, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf("Error calling LLM: %v", err)
	}

	return strings.TrimSpace(string(output))
}
