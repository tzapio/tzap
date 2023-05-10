package embed

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/util"
)

func ProcessDirectory(t *tzap.Tzap, dir string) {
	files := GetGoFilesInDir(dir)
	embeddings, totalTokens, totalLines := ProcessFiles(t, files)
	idsToDelete, err := GetDrift(t, embeddings)
	if err != nil {
		panic(err)
	}
	if err := RemoveOldEmbeddings(t, idsToDelete); err != nil {
		panic(err)
	}
	if err := ProcessEmbeddings(t, embeddings); err != nil {
		panic(err)
	}
	fmt.Printf("Total Files: %d, Total Embeddings: %d, Total Tokens: %d, Total Lines: %d\n", len(files), len(embeddings.Vectors), totalTokens, totalLines)
}
func GetEmbeddings(t *tzap.Tzap, dir string) (types.Embeddings, error) {
	filecontent, err := os.ReadFile("./.tzap-data/files.json")
	if err != nil {
		return types.Embeddings{}, err
	}
	var embeddings types.Embeddings

	if err := json.Unmarshal(filecontent, &embeddings); err != nil {
		return types.Embeddings{}, err
	}
	return embeddings, nil
}
func GetGoFilesInDir(dir string) []string {
	files, err := util.WalkFilesInDir(dir)
	if err != nil {
		panic(err)
	}
	return util.FilterFiles(files, ".go")
}

func ProcessFiles(t *tzap.Tzap, files []string) (types.Embeddings, int, int) {
	totalTokens := 0
	totalLines := 0

	embeddings := types.Embeddings{}

	for _, file := range files {
		fileTokens, lines, content, err := ProcessFile(file, t)
		if err != nil {
			panic(err)
		}
		totalLines += lines
		totalTokens += fileTokens

		//fmt.Printf("File: %s - Tokens: %d, Lines: %d\n", file, fileTokens, lines)

		fileEmbeddings, err := ProcessFileOffsets(t, file, content, fileTokens)
		if err != nil {
			panic(err)
		}

		embeddings.Vectors = append(embeddings.Vectors, fileEmbeddings.Vectors...)
	}

	return embeddings, totalTokens, totalLines
}

func GetDrift(t *tzap.Tzap, nowEmbeddings types.Embeddings) ([]string, error) {
	storedEmbeddings, err := t.TG.ListAllEmbeddingsIds(t.C)
	if err != nil {
		return nil, err
	}

	nowEmbeddingsIds := make(map[string]struct{})
	for _, vectorID := range nowEmbeddings.Vectors {

		nowEmbeddingsIds[vectorID.ID] = struct{}{}
	}
	missingIds := []string{}
	for _, storedVector := range storedEmbeddings.Results {
		if _, exists := nowEmbeddingsIds[storedVector.ID]; !exists {

			missingIds = append(missingIds, storedVector.ID)
			tzap.Log(t, "Drift: ", storedVector)
		}
	}
	println("Drift Check: ", len(storedEmbeddings.Results), len(nowEmbeddings.Vectors), len(missingIds))
	return missingIds, nil
}

func RemoveOldEmbeddings(t *tzap.Tzap, deleteIds []string) error {
	for _, deleteId := range deleteIds {
		if err := t.TG.DeleteEmbeddingDocument(t.C, deleteId); err != nil {
			return err
		}
		tzap.Log(t, "Deleted: ", deleteId)
	}
	return nil
}

func ProcessFileOffsets(t *tzap.Tzap, file, content string, fileTokens int) (types.Embeddings, error) {
	vectors := []types.Vector{}
	start := 0
	end := 250

	for start < fileTokens {
		splitPart, err := ProcessOffset(t, file, content, start, end, fileTokens)
		if err != nil {
			return types.Embeddings{}, err
		}
		vectors = append(vectors, splitPart)

		start += 200
		end += 200
	}

	return types.Embeddings{Vectors: vectors}, nil
}
func ProcessFile(filename string, t *tzap.Tzap) (int, int, string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return 0, 0, "", err
	}

	content := string(fileContent)
	lines := strings.Count(content, "\n")

	fileTokens, err := t.TG.CountTokens(context.Background(), content)
	if err != nil {
		return 0, 0, "", err
	}

	return fileTokens, lines, content, nil
}

func ProcessOffset(t *tzap.Tzap, filename, content string, start int, end int, fileTokens int) (types.Vector, error) {
	truncatedEnd := end
	if end > fileTokens {
		truncatedEnd = fileTokens
	}

	splitPart, err := t.TG.OffsetTokens(t.C, content, start, truncatedEnd)
	if err != nil {
		return types.Vector{}, err
	}

	metadata := map[string]string{
		"filename":     filename,
		"start":        fmt.Sprintf("%d", start),
		"end":          fmt.Sprintf("%d", end),
		"truncatedEnd": fmt.Sprintf("%d", truncatedEnd),
		"splitPart":    splitPart,
		"splitMd5":     util.MD5Hash(splitPart),
	}
	id := metadata["filename"] + "-" + metadata["start"] + "-" + metadata["end"]
	return types.Vector{ID: id, Metadata: metadata}, nil
}