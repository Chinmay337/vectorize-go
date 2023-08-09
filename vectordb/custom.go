package vectordb

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type CollectionParams struct {
	CollectionName     string
	Description        string
	Fields             []*entity.Field
	EnableDynamicField bool
	ShardNum           int32
}

type InsertParams struct {
	CollectionName string
	PartitionName  string
	Columns        map[string]entity.Column
}

func InsertData(milvusClient client.Client, params InsertParams, ctx context.Context) {
	columns := make([]entity.Column, 0, len(params.Columns))
	for _, column := range params.Columns {
		columns = append(columns, column)
	}
	_, err := milvusClient.Insert(
		ctx,                   // ctx
		params.CollectionName, // CollectionName
		params.PartitionName,  // partitionName
		columns...,            // Columns for Collection
	)
	if err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}
	fmt.Printf("Successfully inserted data into %s\n", params.CollectionName)
}

func CreateCollectionFromStruct(milvusClient client.Client, params CollectionParams, ctx context.Context) error {
	schema := &entity.Schema{
		CollectionName:     params.CollectionName,
		Description:        params.Description,
		Fields:             params.Fields,
		EnableDynamicField: params.EnableDynamicField,
	}
	// check if collection exists already
	collections, err := milvusClient.ListCollections(ctx)
	if err != nil {
		log.Fatal("failed to list collections:", err.Error())
	}
	for _, collection := range collections {
		fmt.Println(collection.Name)
		if collection.Name == params.CollectionName {
			fmt.Printf("Collection %s already exists\n", params.CollectionName)
			return nil
		}
	}
	err = milvusClient.CreateCollection(
		ctx, // ctx
		schema,
		params.ShardNum, // shardNum
	)
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}
	fmt.Printf("Successfully created collection %s\n", params.CollectionName)
	return nil
}

// fields := []*entity.Field{
// 	{
// 		Name:     "word",
// 		DataType: entity.FieldTypeVarChar,
// 		TypeParams: map[string]string{
// 			"max_length": "100", // adjust this to the maximum length of your words
// 		},
// 		PrimaryKey: true,
// 		AutoID:     false,
// 	},
// 	{
// 		Name:     "embedding",
// 		DataType: entity.FieldTypeFloatVector,
// 		TypeParams: map[string]string{
// 			"dim": "3", // adjust this to match the dimensionality of your word embeddings
// 		},
// 	},
// }

// params := vectordb.CollectionParams{
// 	CollectionName:     "words",
// 	Description:        "Word embeddings",
// 	Fields:             fields,
// 	EnableDynamicField: true,
// 	ShardNum:           2,
// }

// _ = vectordb.CreateCollectionFromStruct(client, params, ctx)
