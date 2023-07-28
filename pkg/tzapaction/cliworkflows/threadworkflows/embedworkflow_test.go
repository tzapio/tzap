package threadworkflows_test

import (
	"testing"

	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapaction/actionpb"
	"github.com/tzapio/tzap/pkg/tzapaction/cliworkflows/threadworkflows"
)

func TestEmbedWorkflow(t *testing.T) {

	embeddings := []*actionpb.Embedding{
		{
			Content: "Embedding 1",
		},
		{
			Content: "Embedding 2",
		},
	}

	workflow := threadworkflows.EmbedWorkflow(embeddings)
	tzap := tzap.InternalNew()

	result := workflow.Workflow(tzap)
	results := result.GetThread()
	if len(results) != 3 {
		t.Errorf("Expected %d system messages, but got %d", 3, len(results))
	}

	expectedMessages := []string{
		"The following file contents are embeddings for the user input:",
		"Embedding 1",
		"Embedding 2",
	}

	for i, msg := range results {
		if msg.Content != expectedMessages[i] {
			t.Errorf("Expected system message '%s', but got '%s'", expectedMessages[i], msg.Content)
		}
	}
}
