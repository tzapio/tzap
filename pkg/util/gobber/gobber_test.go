package gobber_test

import (
	"bytes"
	"testing"

	"github.com/tzapio/tzap/pkg/util/gobber"
)

func TestGobWriterAndReader(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := gobber.NewGobWriterIO(buf)
	reader := gobber.NewGobReaderIO(buf)

	type Data struct {
		Number int
		Text   string
	}

	writeData := Data{
		Number: 123,
		Text:   "hello world",
	}

	err := writer.Write(writeData)
	if err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	readData := Data{}
	err = reader.Read(&readData)
	if err != nil {
		t.Fatalf("Failed to read data: %v", err)
	}

	if readData != writeData {
		t.Fatalf("Read data does not match written data: got %v, want %v", readData, writeData)
	}
}
func TestGobWriterAndReaderMultiple(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := gobber.NewGobWriterIO(buf)
	reader := gobber.NewGobReaderIO(buf)

	type Data struct {
		Number int
		Text   string
	}

	writeData := []Data{
		{Number: 123, Text: "hello world"},
		{Number: 456, Text: "goodbye world"},
		{Number: 789, Text: "another string"},
	}

	for _, wd := range writeData {
		err := writer.Write(wd)
		if err != nil {
			t.Fatalf("Failed to write data: %v", err)
		}
	}

	for i := 0; i < len(writeData); i++ {
		readData := Data{}
		err := reader.Read(&readData)
		if err != nil {
			t.Fatalf("Failed to read data: %v", err)
		}

		if readData != writeData[i] {
			t.Fatalf("Read data does not match written data: got %v, want %v", readData, writeData[i])
		}
	}
}
