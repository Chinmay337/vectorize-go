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
	fmt.Println("Hello World")
	var ctx context.Context
	ctx = context.Background()
	client := Connect()
	defer client.Close()

	// add timing logs between every function call

	vectorize.Train()
	vectorize.QueryVector("cat")

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
				"dim": "300", // adjust this to match the dimensionality of your word embeddings

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

	vectordb.CreateCollectionFromStruct(client, params, ctx)

	startTime := time.Now()
	vectordb.CreateCollection(client, "books", ctx)
	LogTime(startTime, "CreateCollection")

	startTime = time.Now()
	vectordb.InsertRawVectorIntoCollection(client, "books", ctx)
	LogTime(startTime, "InsertIntoCollection")

	startTime = time.Now()
	vectordb.CreateIndex(client, "books", "book_intro", entity.L2, 1024, ctx)
	LogTime(startTime, "CreateIndex")

	startTime = time.Now()
	vectordb.LoadCollection(client, "books", ctx)
	LogTime(startTime, "LoadCollection")

	startTime = time.Now()
	vectordb.ConductSearch(client, "books", []string{"book_id"}, []float32{0.1, 0.2}, 2, ctx)
	LogTime(startTime, "ConductSearch")

	startTime = time.Now()
	vectordb.DeleteCollection(client, "books", ctx)
	vectordb.DeleteCollection(client, "words", ctx)
	LogTime(startTime, "DeleteCollection")

	// }

	for {
		// Keep App Running - view profile at http://localhost:6060/debug/pprof/
	}
}

func Connect() client.Client {
	milvusClient, err := client.NewGrpcClient( // Max 65,536 connections
		context.Background(), // ctx
		"localhost:19530",    // addr
	)
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

Can also do
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
