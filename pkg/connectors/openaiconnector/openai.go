package openaiconnector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/tiktoken-go/tokenizer"
	"github.com/tzapio/tzap/internal/filelog"
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/connectors/openaiconnector/output"
	"github.com/tzapio/tzap/pkg/types"
)

func InitiateOpenaiClient(apikey string) *OpenaiTgenerator {
	return &OpenaiTgenerator{client: getClient(apikey)}
}

func getClient(apikey string) *openai.Client {
	client := openai.NewClient(apikey)
	return client
}
func (ot OpenaiTgenerator) GenerateChat(ctx context.Context, messages []types.Message) (string, error) {
	config := config.FromContext(ctx)
	content, err := ot.fetchChatResponse(ctx, config.OpenAIModel, false, messages)
	if err != nil {
		filelog.LogData(ctx, err.Error(), filelog.ResponseLog)
		return "", fmt.Errorf("error generating chat prompt result: %v", err)
	}
	return content, nil
}
func (ot OpenaiTgenerator) CountTokens(content string) (int, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		panic("error getting tokenizer")
	}

	ids, _, err := enc.Encode(content)
	if err != nil {
		return 0, errors.New("error couting tokens while encoding")
	}
	return len(ids), err
}

func (ot OpenaiTgenerator) OffsetTokens(content string, from int, to int) (string, error) {
	enc, err := tokenizer.Get(tokenizer.Cl100kBase)
	if err != nil {
		panic("error getting tokenizer")
	}

	ids, _, err := enc.Encode(content)
	if err != nil {
		return "", errors.New("error couting tokens while encoding")
	}
	ids = ids[from:to]
	s, _ := enc.Decode(ids)
	return s, err
}

// fetchChatResponse requests openai-chat completion for the given Tzap and returns the modified content.
func (ot OpenaiTgenerator) fetchChatResponse(ctx context.Context, gptmodel string, stream bool, messages []types.Message) (string, error) {
	// Create a context with a timeout

	request := openai.ChatCompletionRequest{
		Model:       gptmodel,
		Messages:    output.GetOpenAICompletionMessage(messages),
		Temperature: 1,
	}
	filelog.LogData(ctx, request, filelog.RequestLog)
	var content string
	if stream {
		cntnt, err := ot.streamCompletion(request)
		if err != nil {
			return "", fmt.Errorf("chatcompletion error: %v", err)
		}
		content = cntnt
	} else {
		cntnt, err := ot.createChatCompletion(request)
		if err != nil {
			return "", fmt.Errorf("chatcompletion error: %v", err)
		}
		content = cntnt
	}
	filelog.LogData(ctx, content, filelog.ResponseLog)
	return content, nil
}

func (ot OpenaiTgenerator) streamCompletion(request openai.ChatCompletionRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	// Create a stream completion
	s, err := ot.client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		panic(fmt.Errorf("CreateChatCompletionStream Stream error: %v", err))
	}
	var resultBuilder strings.Builder
	// Consume the stream completion
	for {
		// Read the next token from the stream
		response, err := s.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			fmt.Printf("%+v\n", response)
			break
		}
		if err != nil {
			return "", (fmt.Errorf("\n\nStream error: %v", err))
		}
		token := response.Choices[0].Delta.Content
		print(token)
		resultBuilder.WriteString(token)
	}
	response := resultBuilder.String()
	return response, nil
}

func (ot OpenaiTgenerator) createChatCompletion(request openai.ChatCompletionRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	response, err := ot.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", fmt.Errorf("chatcompletion error: %v", err)
	}
	return response.Choices[0].Message.Content, nil
}
