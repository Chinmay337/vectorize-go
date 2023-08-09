package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"milvus/vectordb"
	"milvus/vectorize"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	fmt.Println("Vector DB Golang Interface")
	var ctx context.Context
	ctx = context.Background()
	client := Connect()
	defer client.Close()

	// add timing logs between every function call
	vectorize.Train("string-vectors/input", "string-vectors/word_vector.txt")
	vectorize.QueryVector("cat", "string-vectors/word_vector.txt")

	vectordb.DeleteAllCollections(client, ctx)

	// ====================>>>>>  CREATING COLLECTION

	fields := []*entity.Field{
		{
			Name:     "word",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "100", // adjust this to the maximum length of your words
			},
			PrimaryKey: true,
			AutoID:     false,
		},
		{
			Name:     "embedding",
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "3", // adjust this to match the dimensionality of your word embeddings
			},
		},
	}

	params := vectordb.CollectionParams{
		CollectionName:     "words",
		Description:        "Word embeddings",
		Fields:             fields,
		EnableDynamicField: true,
		ShardNum:           2,
	}

	_ = vectordb.CreateCollectionFromStruct(client, params, ctx)

	// ====================>>>>>  INSERTING into a Collection

	// Assuming using sample embeddings:
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

	// Defining sample Insert Params
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

	// INSERTION DONE ^^^^

	// ====================>>>>>  Searching from a Collection
	// Load Collection
	vectordb.CreateIndex(client, "words", "embedding", entity.L2, 1024, ctx)
	vectordb.LoadCollection(client, "words", ctx)

	fmt.Println("Querying Collection")
	vectordb.QueryCollection(client, "words", "word not in ['cat', 'dog']", []string{"word"}, ctx)

	fmt.Print("\n\nSearching Collection\n\n")

	vectordb.SearchIndexFromCollection(client, "words", "embedding", []entity.Vector{entity.FloatVector([]float32{0.2, 0.2, 0.8})}, []string{"word"}, 3, ctx)

	fmt.Print("\n\nSearching Done\n\n")

	for {
		// Keep App Running - view profile at http://localhost:6060/debug/pprof/
		// Ask user to press enter to exit
		var input string
		fmt.Println("Press Enter to exit")
		fmt.Scanln(&input)
		break
	}
}

func Connect() client.Client {
	var err error
	done := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second)
		select {
		case <-done:
			// If done is closed, it means the function has finished successfully
			return
		default:
			// If done is not closed after 5 seconds, panic
			panic("failed to connect to Milvus. Make sure to run\ndocker-compose up")
		}
	}()

	milvusClient, err := client.NewGrpcClient( // Max 65,536 connections
		context.Background(), // ctx
		"localhost:19530",    // addr
	)
	close(done)
	if err != nil {
		log.Fatal("failed to connect to Milvus:", err.Error())
	}
	fmt.Println("Successfully connected to Milvus")
	return milvusClient
}

func LogTime(startTime time.Time, functionName string) {
	elapsedTime := time.Since(startTime)
	fmt.Printf("Function %s took %s\n\n", functionName, elapsedTime)
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
