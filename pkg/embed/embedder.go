package embed

import (
	"fmt"
	"strings"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
	"github.com/tzapio/tzap/pkg/tzap"
)

type Embedder struct {
	*EmbeddingCache
	*FilestampCache
	EmbedCleaner
}

func NewEmbedder(embeddingCacheDB types.DBCollectionInterface[string], filesTimestampsDB types.DBCollectionInterface[int64]) *Embedder {
	embeddingCache := NewEmbeddingCache(embeddingCacheDB)
	filestampCache := NewFilestampCache(filesTimestampsDB)
	return &Embedder{EmbeddingCache: embeddingCache, EmbedCleaner: EmbedCleaner{}, FilestampCache: filestampCache}
}

func (fe *Embedder) PrepareEmbeddingsFromFiles(t *tzap.Tzap, changedFileContents map[string]string) *types.Embeddings {
	tl.Logger.Println("Preparing embeddings from files", len(changedFileContents))

	rawFileEmbeddings, err := fe.ProcessFileContents(t, changedFileContents)
	if err != nil {
		panic(err)
	}
	return rawFileEmbeddings
}

func (fe *Embedder) ProcessFileContents(t *tzap.Tzap, changedFiles map[string]string) (*types.Embeddings, error) {
	tl.Logger.Println("Processing files", len(changedFiles))
	totalTokens := 0
	totalLines := 0

	embeddings := &types.Embeddings{}
	for file, content := range changedFiles {
		fileTokens, lines, err := fe.ProcessFileContent(t, content)
		if err != nil {
			panic(err)
		}
		totalLines += lines
		totalTokens += fileTokens

		tl.Logger.Printf("File: %s - Tokens: %d, Lines: %d\n", file, fileTokens, lines)

		fileEmbeddings, err := fe.ProcessFileOffsets(t, file, content, fileTokens)
		if err != nil {
			return &types.Embeddings{}, err
		}

		embeddings.Vectors = append(embeddings.Vectors, fileEmbeddings.Vectors...)
	}
	tl.Logger.Println("Processed files", len(changedFiles), "Total Embeddings", len(embeddings.Vectors), "Total Tokens", totalTokens, "Total Lines", totalLines)
	return embeddings, nil
}

func (fe *Embedder) ProcessFileOffsets(t *tzap.Tzap, file string, content string, fileTokens int) (*types.Embeddings, error) {
	vectors := []*types.Vector{}
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
			return &types.Embeddings{}, err
		}

		start := 0
		end := baseStep
		for start < c {
			partialVector, err := fe.ProcessOffset(t, file, chunkContent, start, end, baseStep, chunkStart, lineStart, c)
			if err != nil {
				return &types.Embeddings{}, err
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

	return &types.Embeddings{Vectors: vectors}, nil
}

func (fe *Embedder) ProcessFileContent(t *tzap.Tzap, content string) (int, int, error) {
	lines := strings.Count(content, "\n")

	fileTokens, err := t.TG.CountTokens(t.C, content)
	if err != nil {
		return 0, 0, err
	}

	return fileTokens, lines, nil
}

func (fe *Embedder) ProcessOffset(t *tzap.Tzap, filename, content string, start int, end int, step int, chunkStart int, lineStart int, fileTokens int) (*types.Vector, error) {
	truncatedEnd := end
	if end > fileTokens {
		truncatedEnd = fileTokens
	}
	tl.Logger.Println("Filename", filename, "Processing offset", start, truncatedEnd, "of", fileTokens, "tokens")
	splitPart, _, err := t.TG.OffsetTokens(t.C, content, start, truncatedEnd)
	if err != nil {
		return &types.Vector{}, err
	}
	truncatedRealEnd := start + step
	if truncatedRealEnd > fileTokens {
		truncatedRealEnd = fileTokens
	}

	realSplitPart, _, err := t.TG.OffsetTokens(t.C, content, start, truncatedRealEnd)
	if err != nil {
		return &types.Vector{}, err
	}

	splitPart = "####embedding from file: " + filename + "\n" + splitPart
	metadataStart := chunkStart + start
	metadataEnd := chunkStart + end
	metadataLineStart := lineStart
	metadataTruncatedEnd := truncatedEnd

	id := fmt.Sprintf("%s-%d-%d", filename, metadataStart, metadataEnd)

	return &types.Vector{
		ID: id,
		Metadata: types.Metadata{
			ID:            id,
			Filename:      filename,
			Start:         metadataStart,
			End:           metadataEnd,
			LineStart:     metadataLineStart,
			TruncatedEnd:  metadataTruncatedEnd,
			SplitPart:     splitPart,
			RealSplitPart: realSplitPart,
		},
	}, nil
}
