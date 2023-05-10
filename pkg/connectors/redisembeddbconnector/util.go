package redisembeddbconnector

import (
	"encoding/binary"
	"math"
)

func toBytes(embedding []float32) []byte {
	embeddingBytes := make([]byte, len(embedding)*4)
	for i, f := range embedding {
		binary.LittleEndian.PutUint32(embeddingBytes[i*4:], math.Float32bits(f))
	}
	return embeddingBytes
}
