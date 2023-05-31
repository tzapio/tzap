package localdb_test

import (
	"os"
	"testing"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
)

func TestSet_KeyValueIsSet(t *testing.T) {
	tempFile, err := os.CreateTemp("", "filedb_test")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	db, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		t.Fatalf("Error creating FileDB: %v", err)
	}

	key := "testKey"
	value := "testValue"
	err = db.Set(key, value)
	if err != nil {
		t.Fatalf("Error setting key-value: %v", err)
	}

	readValue, exists := db.Get(key)
	if !exists {
		t.Fatalf("Key %s not found", key)
	}
	if readValue != value {
		t.Fatalf("Value mismatch: expected %s, got %s", value, readValue)
	}
}

func TestBatchSet_NormalValues_MultipleKeyValuesIsSet(t *testing.T) {
	tempFile, err := os.CreateTemp("", "filedb_test")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	db, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		t.Fatalf("Error creating FileDB: %v", err)
	}

	pairs := []types.KeyValue[string]{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "testKey", Value: "value2"},
	}

	wrote, err := db.BatchSet(pairs)
	if err != nil {
		t.Fatalf("Error in BatchSet: %v", err)
	}
	if wrote != len(pairs) {
		t.Fatalf("Expected wrote to be %d, got %d", len(pairs), wrote)
	}
	for _, kv := range pairs {
		readValue, exists := db.Get(kv.Key)
		if !exists {
			t.Fatalf("Key %s not found", kv.Key)
		}
		if readValue != kv.Value {
			t.Fatalf("Value mismatch: expected %s, got %s", kv.Value, readValue)
		}
	}
}

func TestBatchSet_OverridingValues_MultipleKeyValuesIsSet(t *testing.T) {
	tempFile, err := os.CreateTemp("", "filedb_test")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	db, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		t.Fatalf("Error creating FileDB: %v", err)
	}

	pairs := []types.KeyValue[string]{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
		{Key: "testKey", Value: "value2"},
	}

	wrote, err := db.BatchSet(pairs)
	if err != nil {
		t.Fatalf("Error in BatchSet: %v", err)
	}
	if wrote != len(pairs) {
		t.Fatalf("Expected wrote to be %d, got %d", len(pairs), wrote)
	}
	overriding := []types.KeyValue[string]{
		{Key: "key1", Value: "value4"},
		{Key: "key2", Value: "value6"},
		{Key: "testKey", Value: "value28"},
	}
	wrote, err = db.BatchSet(overriding)
	if err != nil {
		t.Fatalf("Error in BatchSet: %v", err)
	}
	if wrote != len(overriding) {
		t.Fatalf("Expected wrote to be %d, got %d", len(overriding), wrote)
	}
	newDbInstance, err := localdb.NewFileDB[string](tempFile.Name())
	if err != nil {
		t.Fatalf("Error creating FileDB: %v", err)
	}
	for _, kv := range overriding {
		readValue, exists := newDbInstance.Get(kv.Key)
		if !exists {
			t.Fatalf("Key %s not found", kv.Key)
		}
		if readValue != kv.Value {
			t.Fatalf("Value mismatch: expected %s, got %s", kv.Value, readValue)
		}
	}
}
