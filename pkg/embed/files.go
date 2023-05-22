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

func LazyPrepareEmbeddingsFromFiles(t *tzap.Tzap, files []string) types.Embeddings {
	rawFileEmbeddings, totalTokens, totalLines := ProcessFiles(t, files)
	idsToDelete, err := GetDrift(t, rawFileEmbeddings)
	if err != nil {
		panic(err)
	}
	if err := RemoveOldEmbeddings(t, idsToDelete); err != nil {
		panic(err)
	}
	uncachedEmbeddings := GetUncachedEmbeddings(rawFileEmbeddings)
	if len(uncachedEmbeddings.Vectors) > 0 {
		if err := FetchAndCacheNewEmbeddings(t, uncachedEmbeddings); err != nil {
			panic(err)
		}
	}

	cachedEmbeddings := GetCachedEmbeddings(rawFileEmbeddings)
	SaveEmbeddingToFile(cachedEmbeddings)

	fmt.Printf("Total Files: %d, Total Embeddings: %d, Total Tokens: %d, Total Lines: %d\n", len(files), len(rawFileEmbeddings.Vectors), totalTokens, totalLines)
	return cachedEmbeddings
}

func PrepareEmbeddingsFromFiles(t *tzap.Tzap, files []string) types.Embeddings {
	rawFileEmbeddings, _, _ := ProcessFiles(t, files)
	idsToDelete, err := GetDrift(t, rawFileEmbeddings)
	if err != nil {
		panic(err)
	}
	if err := RemoveOldEmbeddings(t, idsToDelete); err != nil {
		panic(err)
	}

	//fmt.Printf("Total Files: %d, Total Embeddings: %d, Total Tokens: %d, Total Lines: %d\n", len(files), len(rawFileEmbeddings.Vectors), totalTokens, totalLines)
	return rawFileEmbeddings
}

func GetEmbeddingsFromFile() (types.Embeddings, error) {
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
	//println("Drift Check: ", len(storedEmbeddings.Results), len(nowEmbeddings.Vectors), len(missingIds))
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
	step := 200
	start := 0
	offset := 50
	end := step + offset

	lineStart := 1
	for start < fileTokens {
		partialVector, err := ProcessOffset(t, file, content, start, end, step, lineStart, fileTokens)
		if err != nil {
			return types.Embeddings{}, err
		}
		vectors = append(vectors, partialVector)

		start += step
		end += step
		lines := strings.Count(partialVector.Metadata["realSplitPart"], "\n")
		lineStart += lines
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

func ProcessOffset(t *tzap.Tzap, filename, content string, start int, end int, step int, lineStart int, fileTokens int) (types.Vector, error) {
	truncatedEnd := end
	if end > fileTokens {
		truncatedEnd = fileTokens
	}

	splitPart, err := t.TG.OffsetTokens(t.C, content, start, truncatedEnd)
	if err != nil {
		return types.Vector{}, err
	}
	truncatedRealEnd := start + step
	if truncatedRealEnd > fileTokens {
		truncatedRealEnd = fileTokens
	}
	realSplitPart, err := t.TG.OffsetTokens(t.C, content, start, truncatedRealEnd)
	if err != nil {
		return types.Vector{}, err
	}

	splitPart = "###embedding from file: " + filename + "\n" + splitPart
	metadataStart := fmt.Sprintf("%d", start)
	metadataEnd := fmt.Sprintf("%d", end)
	metadataLineStart := fmt.Sprintf("%d", lineStart)
	metadataTruncatedEnd := fmt.Sprintf("%d", truncatedEnd)
	metadataSplitMd5 := util.MD5Hash(splitPart)

	id := filename + "-" + metadataStart + "-" + metadataEnd

	metadata := map[string]string{
		"id":            id,
		"filename":      filename,
		"start":         metadataStart,
		"end":           metadataEnd,
		"lineStart":     metadataLineStart,
		"truncatedEnd":  metadataTruncatedEnd,
		"splitPart":     splitPart,
		"realSplitPart": realSplitPart,
		"splitMd5":      metadataSplitMd5,
	}
	return types.Vector{ID: id, Metadata: metadata}, nil
}
