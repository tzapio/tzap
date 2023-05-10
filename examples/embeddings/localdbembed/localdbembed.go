package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/tzapio/tzap/pkg/connectors/localdbconnector"
	"github.com/tzapio/tzap/pkg/types"
)

func main() {
	index, err := localdbconnector.InitiateLocalDB("./.tzap-data/fileembeddings.db")
	if err != nil {
		panic(err)
	}

	filecontent, err := os.ReadFile("./.tzap-data/files.json")
	if err != nil {
		panic(err)
	}
	var data types.Embeddings

	if err := json.Unmarshal(filecontent, &data); err != nil {
		panic(err)
	}
	for _, vector := range data.Vectors {
		err := index.AddEmbeddingDocument(context.Background(), vector.ID, vector.Values, vector.Metadata)
		if err != nil {
			println(err.Error())
		}
	}
	println("adding", len(data.Vectors), "documents")
	queryContent, err := os.ReadFile("query.json")
	if err != nil {
		panic(err)
	}
	var query types.QueryJson
	if err := json.Unmarshal(queryContent, &query); err != nil {
		panic(err)
	}

	results, err := index.SearchWithEmbedding(context.Background(), types.QueryFilter{Values: query.Queries[0].Values}, 10)
	if err != nil {
		panic(err)
	}
	for _, result := range results.Results {
		println(result.ID, result.Metadata["__vec_score"], result.Metadata["filename"])
	}
	/*
		println("results", )
		println("results", len(result))
		for _, r := range result {
			println(r.Id, r.Properties["__vec_score"].(string), r.Properties["filename"].(string), r.Properties["splitPart"].(string))
		}*/

}
