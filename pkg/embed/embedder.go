package embed

import (
	"context"
	"fmt"
	"strings"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

type Embedder struct {
	t                 *tzap.Tzap
	filesTimestampsDB *localdb.FileDB[int64]
	*EmbeddingCache
}

func NewEmbedder(t *tzap.Tzap) *Embedder {
	filesTimestampsDB, err := localdb.NewFileDB[int64]("./.tzap-data/filesTimestamps.db")
	if err != nil {
		panic(err)
	}
	embeddingCache := NewEmbeddingCache(filesTimestampsDB)
	return &Embedder{t: t, filesTimestampsDB: filesTimestampsDB, EmbeddingCache: embeddingCache}
}

func (fe *Embedder) PrepareEmbeddingsFromFiles(files []types.FileReader) types.Embeddings {
	tl.Logger.Println("Preparing embeddings from files", len(files))
	changedFiles, unchangedFiles, err := fe.CheckFileCache(files)
	if err != nil {
		panic(err)
	}
	rawFileEmbeddings, _, _ := fe.ProcessFiles(changedFiles)
	storedEmbeddings, err := fe.t.TG.ListAllEmbeddingsIds(fe.t.C)
	if err != nil {
		panic(err)
	}

	idsToDelete, err := fe.GetDrift(storedEmbeddings, rawFileEmbeddings, unchangedFiles)
	if err != nil {
		panic(err)
	}
	tl.Logger.Println("Removing old embeddings", len(idsToDelete))
	if err := fe.RemoveOldEmbeddings(idsToDelete); err != nil {
		panic(err)
	}
	return rawFileEmbeddings
}

func (fe *Embedder) ProcessFiles(changedFiles map[string]string) (types.Embeddings, int, int) {
	tl.Logger.Println("Processing files", len(changedFiles))
	totalTokens := 0
	totalLines := 0

	embeddings := types.Embeddings{}
	for file, content := range changedFiles {
		fileTokens, lines, err := fe.ProcessFile(content)
		if err != nil {
			panic(err)
		}
		totalLines += lines
		totalTokens += fileTokens

		//fmt.Printf("File: %s - Tokens: %d, Lines: %d\n", file, fileTokens, lines)

		fileEmbeddings, err := fe.ProcessFileOffsets(file, content, fileTokens)
		if err != nil {
			panic(err)
		}

		embeddings.Vectors = append(embeddings.Vectors, fileEmbeddings.Vectors...)
	}
	tl.Logger.Println("Processed files", len(changedFiles), "Total Embeddings", len(embeddings.Vectors), "Total Tokens", totalTokens, "Total Lines", totalLines)
	return embeddings, totalTokens, totalLines
}

func (fe *Embedder) GetDrift(storedEmbeddings types.SearchResults, nowEmbeddings types.Embeddings, unchangedFiles map[string]int64) ([]string, error) {
	nowEmbeddingsIds := make(map[string]struct{})
	for _, vectorID := range nowEmbeddings.Vectors {
		nowEmbeddingsIds[vectorID.ID] = struct{}{}
	}
	missingIds := []string{}
	for _, storedVector := range storedEmbeddings.Results {
		filename := storedVector.Vector.Metadata.Filename
		if _, exists := unchangedFiles[filename]; exists {
			continue
		}
		if _, exists := nowEmbeddingsIds[storedVector.Vector.ID]; !exists {

			missingIds = append(missingIds, storedVector.Vector.ID)
			tl.Logger.Println("Drift: ", storedVector)
		}
	}
	//println("Drift Check: ", len(storedEmbeddings.Results), len(nowEmbeddings.Vectors), len(missingIds))
	return missingIds, nil
}
func (fe *Embedder) RemoveOldEmbeddings(deleteIds []string) error {
	if err := fe.t.TG.DeleteEmbeddingDocuments(fe.t.C, deleteIds); err != nil {
		return err
	}
	return nil
}
func (fe *Embedder) ProcessFileOffsets(file string, content string, fileTokens int) (types.Embeddings, error) {
	vectors := []types.Vector{}
	baseStep := 200
	step := 4000
	chunkStart := 0
	chunkEnd := step

	lineStart := 1
	for chunkStart < fileTokens {
		// Process file in chunks of 32k tokens to avoid loading whole file
		tl.Logger.Println("Processing file", file, "chunk", chunkStart, "to", chunkEnd, "tokens", fileTokens)
		chunkContent, c, err := fe.t.TG.OffsetTokens(fe.t.C, content, chunkStart, chunkEnd)
		if err != nil {
			return types.Embeddings{}, err
		}

		start := 0
		end := baseStep
		for start < c {
			partialVector, err := fe.ProcessOffset(file, chunkContent, start, end, baseStep, chunkStart, lineStart, c)
			if err != nil {
				return types.Embeddings{}, err
			}
			vectors = append(vectors, partialVector)

			start += baseStep
			end += baseStep

			lines := strings.Count(partialVector.Metadata.RealSplitPart, "\n")
			lineStart += lines
		}

		chunkStart += step
		chunkEnd += step
	}

	return types.Embeddings{Vectors: vectors}, nil
}
func (fe *Embedder) ProcessFile(content string) (int, int, error) {
	lines := strings.Count(content, "\n")

	fileTokens, err := fe.t.TG.CountTokens(context.Background(), content)
	if err != nil {
		return 0, 0, err
	}

	return fileTokens, lines, nil
}

func (fe *Embedder) ProcessOffset(filename, content string, start int, end int, step int, chunkStart int, lineStart int, fileTokens int) (types.Vector, error) {
	truncatedEnd := end
	if end > fileTokens {
		truncatedEnd = fileTokens
	}
	tl.Logger.Println("Filename", filename, "Processing offset", start, truncatedEnd, "of", fileTokens, "tokens")
	splitPart, _, err := fe.t.TG.OffsetTokens(fe.t.C, content, start, truncatedEnd)
	if err != nil {
		return types.Vector{}, err
	}
	truncatedRealEnd := start + step
	if truncatedRealEnd > fileTokens {
		truncatedRealEnd = fileTokens
	}

	realSplitPart, _, err := fe.t.TG.OffsetTokens(fe.t.C, content, start, truncatedRealEnd)
	if err != nil {
		return types.Vector{}, err
	}

	splitPart = "####embedding from file: " + filename + "\n" + splitPart
	metadataStart := chunkStart + start
	metadataEnd := chunkStart + end
	metadataLineStart := lineStart
	metadataTruncatedEnd := truncatedEnd

	id := fmt.Sprintf("%s-%d-%d", filename, metadataStart, metadataEnd)

	return types.Vector{ID: id, Metadata: types.Metadata{
		ID:            id,
		Filename:      filename,
		Start:         metadataStart,
		End:           metadataEnd,
		LineStart:     metadataLineStart,
		TruncatedEnd:  metadataTruncatedEnd,
		SplitPart:     splitPart,
		RealSplitPart: realSplitPart,
	}}, nil
}
