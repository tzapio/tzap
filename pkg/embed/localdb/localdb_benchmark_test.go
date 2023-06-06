package localdb_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
)

func BenchmarkFileDB_Set(b *testing.B) {
	tempFile, err := os.CreateTemp("", "filedb_test")
	if err != nil {
		b.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	db, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		b.Fatalf("Error creating FileDB: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		err = db.Set(key, value)
		if err != nil {
			b.Fatalf("Error setting key-value: %v", err)
		}
	}
}
func BenchmarkFileDB_BatchSet(b *testing.B) {
	tempFile, err := os.CreateTemp("", "filedb_test")
	if err != nil {
		b.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	db, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		b.Fatalf("Error creating FileDB: %v", err)
	}

	pairs := make([]types.KeyValue[string], b.N)
	for i := 0; i < b.N; i++ {
		pairs[i] = types.KeyValue[string]{
			Key:   fmt.Sprintf("key%d", i),
			Value: fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()
	wrote, err := db.BatchSet(pairs)
	if err != nil {
		b.Fatalf("Error in BatchSet: %v", err)
	}
	if wrote != b.N {
		b.Fatalf("Expected wrote to be %d, but got %d", b.N, wrote)
	}
}
func BenchmarkFileDB_InMemory_BatchSet(b *testing.B) {
	db, err := localdb.NewFileDB[string]("@MEMORY/TEST")
	if err != nil {
		b.Fatalf("Error creating FileDB: %v", err)
	}

	pairs := make([]types.KeyValue[string], b.N)
	for i := 0; i < b.N; i++ {
		pairs[i] = types.KeyValue[string]{
			Key:   fmt.Sprintf("key%d", i),
			Value: fmt.Sprintf("value%d", i),
		}
	}

	b.ResetTimer()
	wrote, err := db.BatchSet(pairs)
	if err != nil {
		b.Fatalf("Error in BatchSet: %v", err)
	}
	if wrote != b.N {
		b.Fatalf("Expected wrote to be %d, but got %d", b.N, wrote)
	}
}
