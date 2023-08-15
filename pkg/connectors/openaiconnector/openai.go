package openaiconnector

import (
	"context"
	"encoding/json"
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

func InitiateOpenaiClient(apikey string, conf config.Configuration) *OpenaiTgenerator {
	tl.Logger.Println("Initiating OpenAI Client")
	tokenizer := tokenizer.NewTokenizer()

	return &OpenaiTgenerator{completionClient: getClient(conf.CompletionURL, apikey), embeddingClient: getClient(conf.EmbeddingURL, apikey), Tokenizer: tokenizer}
}

func getClient(baseurl, apikey string) *openai.Client {
	if baseurl == "" {
		return openai.NewClient(apikey)
	}

	config := openai.DefaultConfig(apikey)
	config.BaseURL = baseurl
	client := openai.NewClientWithConfig(config)

	return client
}

func (ot *OpenaiTgenerator) GenerateChat(ctx context.Context, messages []types.Message, stream bool, functions string) (types.CompletionMessage, error) {
	config := config.FromContext(ctx)
	content, err := ot.fetchChatResponse(ctx, config.OpenAIModel, stream, messages, functions)
	if err != nil {
		return types.CompletionMessage{}, fmt.Errorf("error generating chat prompt result: %v", err)
	}
	return content, nil
}

// fetchChatResponse requests openai-chat completion for the given Tzap and returns the modified content.
func (ot *OpenaiTgenerator) fetchChatResponse(ctx context.Context, gptmodel string, stream bool, messages []types.Message, functions string) (types.CompletionMessage, error) {
	// Create a context with a timeout
	config := config.FromContext(ctx)
	var functionDefinitions []openai.FunctionDefinition
	if functions != "" {
		err := json.Unmarshal([]byte(functions), &functionDefinitions)
		if err != nil {
			return types.CompletionMessage{}, fmt.Errorf("error unmarshalling function definitions: %v", err)
		}
	}

	request := openai.ChatCompletionRequest{
		Model:       gptmodel,
		Messages:    output.GetOpenAICompletionMessage(messages),
		Temperature: config.Temperature,
		Functions:   functionDefinitions,
	}

	if stream {
		streamContent, err := ot.streamCompletion(ctx, request)
		if err != nil {
			return types.CompletionMessage{}, fmt.Errorf("chatcompletion error: %v", err)
		}
		return streamContent, nil
	} else {
		responseContent, err := ot.createChatCompletion(ctx, request)
		if err != nil {
			return types.CompletionMessage{}, fmt.Errorf("chatcompletion error: %v", err)
		}
		return responseContent, nil
	}

}

func (ot *OpenaiTgenerator) streamCompletion(ctx context.Context, request openai.ChatCompletionRequest) (types.CompletionMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	// Create a stream completion
	retries := 3

	for i := 0; i < retries; i++ {
		s, err := ot.completionClient.CreateChatCompletionStream(ctx, request)
		if err != nil {
			e := &openai.APIError{}
			if errors.As(err, &e) {
				switch e.HTTPStatusCode {
				case 401:
					panic(fmt.Errorf("invalid openai key. Please check your key and try again. %v", err))
				case 429:
					// rate limiting or engine overload (wait and retry)
					tl.Logger.Println("rate limiting or engine overload (wait and retry)")
					continue
				case 500:
					// openai server error (retry)
					tl.Logger.Println("openai server error (retry)")
					continue
				default:
					return types.CompletionMessage{}, fmt.Errorf("stream error: %v", err)
				}
			}
		}

		var resultBuilder strings.Builder
		// Consume the stream completion
		var functionCall types.FunctionCall
		var finishReason types.FinishReason
		for {
			// Read the next token from the stream
			response, err := s.Recv()
			if errors.Is(err, io.EOF) {
				//	fmt.Println("\nStream finished")
				//	fmt.Printf("%+v\n", response)
				break
			}
			if err != nil {
				finishReason = "error"
				break
			}

			if response.Choices[0].FinishReason != "" {
				finishReason = types.FinishReason(response.Choices[0].FinishReason)
			}
			if response.Choices[0].Delta.FunctionCall != nil {
				q := response.Choices[0].Delta.FunctionCall
				if functionCall.Name != q.Name && q.Name != "" {
					functionCall.Name = q.Name
					println(q.Name)
				}
				print(q.Arguments)
				resultBuilder.WriteString(q.Arguments)
			} else {
				token := response.Choices[0].Delta.Content
				print(token)
				resultBuilder.WriteString(token)
			}
		}
		if functionCall.Name != "" {
			functionCall.Arguments = resultBuilder.String()
			return types.CompletionMessage{FunctionCall: &functionCall, FinishReason: finishReason}, nil
		} else {
			return types.CompletionMessage{Content: resultBuilder.String(), FinishReason: finishReason}, nil
		}
	}
	return types.CompletionMessage{}, errors.New("stream error: retries exceeded")
}

func (ot *OpenaiTgenerator) createChatCompletion(ctx context.Context, request openai.ChatCompletionRequest) (types.CompletionMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	response, err := ot.completionClient.CreateChatCompletion(ctx, request)
	if err != nil {
		return types.CompletionMessage{}, fmt.Errorf("chatcompletion error: %v", err)
	}
	if response.Choices[0].Message.FunctionCall.Name != "" {
		return types.CompletionMessage{FinishReason: types.FinishReason(response.Choices[0].FinishReason), FunctionCall: &types.FunctionCall{Name: response.Choices[0].Message.FunctionCall.Name, Arguments: response.Choices[0].Message.FunctionCall.Arguments}}, nil
	}
	return types.CompletionMessage{FinishReason: types.FinishReason(response.Choices[0].FinishReason), Content: response.Choices[0].Message.Content}, nil
}
