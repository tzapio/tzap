package openaiconnector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector/output"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector/tokenizer"

	"github.com/tzapio/tzap/pkg/types"
)

func InitiateOpenaiClient(apikey string) *OpenaiTgenerator {
	tl.Logger.Println("Initiating OpenAI Client")
	tokenizer := tokenizer.NewTokenizer()

	return &OpenaiTgenerator{client: getClient(apikey), Tokenizer: tokenizer}
}

func getClient(apikey string) *openai.Client {
	client := openai.NewClient(apikey)
	return client
}
func (ot *OpenaiTgenerator) GenerateChat(ctx context.Context, messages []types.Message, stream bool) (string, error) {
	config := config.FromContext(ctx)
	content, err := ot.fetchChatResponse(ctx, config.OpenAIModel, stream, messages)
	if err != nil {
		return "", fmt.Errorf("error generating chat prompt result: %v", err)
	}
	return content, nil
}

// fetchChatResponse requests openai-chat completion for the given Tzap and returns the modified content.
func (ot *OpenaiTgenerator) fetchChatResponse(ctx context.Context, gptmodel string, stream bool, messages []types.Message) (string, error) {
	// Create a context with a timeout
	config := config.FromContext(ctx)
	request := openai.ChatCompletionRequest{
		Model:       gptmodel,
		Messages:    output.GetOpenAICompletionMessage(messages),
		Temperature: config.Temperature,
	}
	var content string
	if stream {
		streamContent, err := ot.streamCompletion(ctx, request)
		if err != nil {
			return "", fmt.Errorf("chatcompletion error: %v", err)
		}
		content = streamContent
	} else {
		responseContent, err := ot.createChatCompletion(ctx, request)
		if err != nil {
			return "", fmt.Errorf("chatcompletion error: %v", err)
		}
		content = responseContent
	}
	return content, nil
}

func (ot *OpenaiTgenerator) streamCompletion(ctx context.Context, request openai.ChatCompletionRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	// Create a stream completion
	retries := 3
	for i := 0; i < retries; i++ {
		s, err := ot.client.CreateChatCompletionStream(ctx, request)
		if err != nil {
			e := &openai.APIError{}
			if errors.As(err, &e) {
				switch e.HTTPStatusCode {
				case 401:
					panic(fmt.Errorf("invalid openai key. Please check your key and try again. %v", err))
				case 429:
					// rate limiting or engine overload (wait and retry)
					continue
				case 500:
					// openai server error (retry)
					continue
				default:
					return "", fmt.Errorf("stream error: %v", err)
				}
			}
		}

		var resultBuilder strings.Builder
		// Consume the stream completion

		for {
			// Read the next token from the stream
			response, err := s.Recv()
			if errors.Is(err, io.EOF) {
				//	fmt.Println("\nStream finished")
				//	fmt.Printf("%+v\n", response)
				break
			}
			if err != nil {
				return resultBuilder.String(), fmt.Errorf("stream error: %v", err)
			}

			token := response.Choices[0].Delta.Content
			print(token)
			resultBuilder.WriteString(token)
		}
		response := resultBuilder.String()
		return response, nil
	}
	return "", errors.New("stream error: retries exceeded")
}

func (ot *OpenaiTgenerator) createChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	response, err := ot.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", fmt.Errorf("chatcompletion error: %v", err)
	}
	return response.Choices[0].Message.Content, nil
}
