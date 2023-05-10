// file: pkg/embed/filesbatch_test.go
package embed_test

import (
	"testing"
)

func Test_SplitCachedUncachedEmbeddings_GivenEmbeddings_SplitsCorrectly(t *testing.T) {
	/*
		TODO requires refactor to automate


		// Create a temporary test cache file
		cacheFile, err := ioutil.TempFile("", "testembeddingsCache.db")
		require.NoError(t, err)
		defer os.Remove(cacheFile.Name())

		// Create input embeddings
		testEmbeddings := types.Embeddings{
			Vectors: []types.Vector{
				{
					ID: "test_vector_1",
					Metadata: map[string]string{
						"splitPart": "test_key",
					},
				},
				{
					ID: "test_vector_2",
					Metadata: map[string]string{
						"splitPart": "non_cached_key",
					},
				},
			},
		}

		// Perform the SplitCachedUncachedEmbeddings function
		cached, uncached := embed.SplitCachedUncachedEmbeddings(testEmbeddings)

		// Check if the result is as expected
		expectedCached := []types.Vector{
			{
				ID: "test_vector_1",
				Metadata: map[string]string{
					"splitPart": "test_key",
				},
				Values: []float32{},
			},
		}

		expectedUncached := []types.Vector{
			{
				ID: "test_vector_2",
				Metadata: map[string]string{
					"splitPart": "non_cached_key",
				},
			},
		}

		assert.Equal(t, expectedCached, cached)
		assert.Equal(t, expectedUncached, uncached)*/
}

func Test_SaveBatchToFile_GivenBatch_SavesToFile(t *testing.T) {
	/*
		TODO requires refactor to automate
		// Create test input data
		testBatch := []types.Vector{
			{
				ID: "test_vector_1",
				Metadata: map[string]string{
					"batchNumber": "1",
				},
				Values: []float32{1.2, 2.3, 3.4},
			},
		}

		// Perform the SaveBatchToFile function
		err := embed.SaveBatchToFile(testBatch, 1)
		if err != nil {
			t.Errorf("Error saving batch to file: %s", err)
		}

		// Check if the file is created and the content is correct
		content, err := os.ReadFile("files-1.json")
		if err != nil {
			t.Errorf("Error reading file: %s", err)
		}
		os.Remove("files-1.json")

		var parsedContent types.Embeddings
		err = json.Unmarshal(content, &parsedContent)
		if err != nil {
			t.Errorf("Error parsing json: %s", err)
		}
		if len(testBatch) != len(parsedContent.Vectors) {
			t.Errorf("Expected %v, got %v", testBatch, parsedContent.Vectors)
		}*/
}

func Test_BatchEmbeddings_GivenEmbeddings_BatchesCorrectly(t *testing.T) {
	/*
		TODO requires refactor to automate


		// Create test input data
		testEmbeddings := types.Embeddings{
			Vectors: []types.Vector{
				{
					ID: "test_vector_1",
					Metadata: map[string]string{
						"batchNumber": "1",
					},
					Values: []float32{1.0, 2.0, 3.0},
				},
				{
					ID: "test_vector_2",
					Metadata: map[string]string{
						"batchNumber": "2",
					},
					Values: []float32{4.0, 5.0, 6.0},
				},
			},
		}

		// Perform the BatchEmbeddings function
		err := embed.BatchEmbeddings(testEmbeddings)
		require.NoError(t, err)

		// Check if the files are created and the content is correct
		for i := 1; i <= 2; i++ {
			filename := fmt.Sprintf("files-%d.json", i)
			content, err := ioutil.ReadFile(filename)
			require.NoError(t, err)
			os.Remove(filename)

			var parsedContent types.Embeddings
			err = json.Unmarshal(content, &parsedContent)
			require.NoError(t, err)

			assert.Equal(t, testEmbeddings.Vectors[i-1:i], parsedContent.Vectors)
		}*/
}
