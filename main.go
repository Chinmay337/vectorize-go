package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"

	"milvus/tools"
	"milvus/vectordb"
	"milvus/vectorize"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func main() {
	tools.EnablePerformanceServerIfFlag()

	var ctx context.Context
	ctx = context.Background()

	client := tools.ConnectVectorDB()
	defer client.Close()

	// ----->>  Converting Raw UTF-8 Strings to Embeddings from a raw UTF-8 file  <<--------

	vectorize.Train("string-vectors/input", "string-vectors/word_vector.txt")
	vectorize.QueryVector("cat", "string-vectors/word_vector.txt")

	vectordb.DeleteAllCollections(client, ctx)

	// ------------>  CREATING COLLECTIONS  <------------
	_ = vectordb.NewCollectionBuilder().
		WithName("words").
		WithDescription("collection of words").
		WithFields(
			vectordb.NewFieldVarChar("word", 100, true, false),
			vectordb.NewFieldFloatVector("embedding", 3),
		).
		Create(client, ctx)

	// ------------>  INSERTING into a Collection   <------------

	// Example using Raw Embeddings

	words := []string{"word1", "word2", "word3", "cat", "dog"}
	embeddings := [][]float32{
		{0.1, 0.2, 0.3}, // embedding for word1
		{0.4, 0.5, 0.6}, // embedding for word2
		{0.7, 0.8, 0.9}, // embedding for word3
		{0.2, 0.2, 0.7}, // embedding for cat
		{0.2, 0.2, 0.8}, // embedding for dog
	}

	wordColumn := entity.NewColumnVarChar("word", words)
	embeddingColumn := entity.NewColumnFloatVector("embedding", 3, embeddings) // 3 is the dimensionality of the embeddings

	// Defining Insert Params
	insertParams := vectordb.InsertParams{
		CollectionName: "words",
		PartitionName:  "", // specify partition name if needed
		Columns: map[string]entity.Column{
			"word":      wordColumn,
			"embedding": embeddingColumn,
		},
	}

	// Inserting into the Collection
	vectordb.InsertData(client, insertParams, ctx)

	// ------------>  Searching from a Collection  <------------

	vectordb.CreateIndex(client, "words", "embedding", entity.L2, 1024, ctx)

	vectordb.LoadCollection(client, "words", ctx)

	vectordb.QueryCollection(client, "words", "word not in ['cat', 'dog']", []string{"word"}, ctx)

	fmt.Print("\n\nSearching Collection\n\n")

	vectordb.SearchIndexFromCollection(client, "words", "embedding", []entity.Vector{entity.FloatVector([]float32{0.2, 0.2, 0.8})}, []string{"word"}, 3, ctx)

	fmt.Print("\n\nSearching Done\n\n")

	for {
		// Keep App Running - view profile at http://localhost:6060/debug/pprof/
		var input string
		log.Println("Press Enter to exit")
		fmt.Scanln(&input)
		break
	}
}

/*

BENCHMARKING

go build -gcflags '-m -l' main.go

go-torch -u http://localhost:6060/

go tool pprof -seconds 30 http://localhost:6060/debug/pprof/profile

For Web UI Viz
go tool pprof -http=localhost:8080 http://localhost:6060/debug/pprof/profile

For Memory Profile with Web UI Viz
go tool pprof -http=localhost:8080 http://localhost:6060/debug/pprof/heap

inuse_space:   memory allocated but not yet released
inuse_objects:  objects allocated but not yet released

alloc_space:   the  total amount of memory allocated
alloc_objects : the  total number of objects allocated

*/
