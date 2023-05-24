package embed

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tzapio/tzap/internal/logging/tl"
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
	tl.Logger.Println("Preparing embeddings from files", len(files))
	rawFileEmbeddings, _, _ := ProcessFiles(t, files)
	idsToDelete, err := GetDrift(t, rawFileEmbeddings)
	if err != nil {
		panic(err)
	}
	tl.Logger.Println("Removing old embeddings", len(idsToDelete))
	if err := RemoveOldEmbeddings(t, idsToDelete); err != nil {
		panic(err)
	}
	return rawFileEmbeddings
}

func GetEmbeddingsFromFile() (types.Embeddings, error) {
	tl.Logger.Println("Getting embeddings from file", "./.tzap-data/files.json")
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
	tl.Logger.Println("Processing files", len(files))
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
	tl.Logger.Println("Processed files", len(files), "Total Embeddings", len(embeddings.Vectors), "Total Tokens", totalTokens, "Total Lines", totalLines)
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
func ProcessFileOffsets(t *tzap.Tzap, file string, content string, fileTokens int) (types.Embeddings, error) {
	vectors := []types.Vector{}
	baseStep := 200
	step := 4000
	chunkStart := 0
	chunkEnd := step

	lineStart := 1
	for chunkStart < fileTokens {
		// Process file in chunks of 32k tokens to avoid loading whole file
		tl.Logger.Println("Processing file", file, "chunk", chunkStart, "to", chunkEnd, "tokens", fileTokens)
		chunkContent, c, err := t.TG.OffsetTokens(t.C, content, chunkStart, chunkEnd)
		if err != nil {
			return types.Embeddings{}, err
		}

		start := 0
		end := baseStep
		for start < c {
			partialVector, err := ProcessOffset(t, file, chunkContent, start, end, baseStep, chunkStart, lineStart, c)
			if err != nil {
				return types.Embeddings{}, err
			}
			vectors = append(vectors, partialVector)

			start += baseStep
			end += baseStep

			lines := strings.Count(partialVector.Metadata["realSplitPart"], "\n")
			lineStart += lines
		}

		chunkStart += step
		chunkEnd += step
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

func ProcessOffset(t *tzap.Tzap, filename, content string, start int, end int, step int, chunkStart int, lineStart int, fileTokens int) (types.Vector, error) {
	truncatedEnd := end
	if end > fileTokens {
		truncatedEnd = fileTokens
	}
	tl.Logger.Println("Filename", filename, "Processing offset", start, truncatedEnd, "of", fileTokens, "tokens")
	splitPart, _, err := t.TG.OffsetTokens(t.C, content, start, truncatedEnd)
	if err != nil {
		return types.Vector{}, err
	}
	truncatedRealEnd := start + step
	if truncatedRealEnd > fileTokens {
		truncatedRealEnd = fileTokens
	}
	realSplitPart, _, err := t.TG.OffsetTokens(t.C, content, start, truncatedRealEnd)
	if err != nil {
		return types.Vector{}, err
	}

	splitPart = "###embedding from file: " + filename + "\n" + splitPart
	metadataStart := fmt.Sprintf("%d", chunkStart+start)
	metadataEnd := fmt.Sprintf("%d", chunkStart+end)
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
