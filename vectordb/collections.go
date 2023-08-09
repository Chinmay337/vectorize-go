package vectordb

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)



func CreateCollection(milvusClient client.Client, collection string, ctx context.Context) {
	schema := &entity.Schema{
		CollectionName: collection,
		Description:    "Test book search",
		Fields: []*entity.Field{
			{
				Name:       "book_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:       "word_count",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: false,
				AutoID:     false,
			},
			{
				Name:     "book_intro",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "2",
				},
			},
		},
		EnableDynamicField: true,
	}
	err := milvusClient.CreateCollection(
		ctx, // ctx
		schema,
		2, // shardNum
	)
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}
	fmt.Printf("Successfully created collection %s\n", collection)
}

func DeleteCollection(milvusClient client.Client, collection string, ctx context.Context) error {
	err := milvusClient.DropCollection(
		context.Background(), // ctx
		collection,           // CollectionName
	)
	if err != nil {
		log.Fatal("fail to drop collection:", err.Error())
	}
	fmt.Println("Successfully dropped collection:", collection)
	return nil
}

func DeleteAllCollections(milvusClient client.Client, ctx context.Context) error {
	collections, err := milvusClient.ListCollections(ctx)
	if err != nil {
		log.Fatal("failed to list collections:", err.Error())
	}
	for _, collection := range collections {
		err = milvusClient.DropCollection(
			ctx, // ctx
			collection.Name, // CollectionName
		)
		if err != nil {
			log.Fatal("fail to drop collection:", err.Error())
		}
		fmt.Println("Successfully dropped collection:", collection.Name)
	}
	return nil
}

func InsertRawVectorIntoCollection(milvusClient client.Client, collection string, ctx context.Context) {
	// Prepare Data
	bookIDs := make([]int64, 0, 2000)
	wordCounts := make([]int64, 0, 2000)
	bookIntros := make([][]float32, 0, 2000)
	for i := 0; i < 2000; i++ {
		bookIDs = append(bookIDs, int64(i))
		wordCounts = append(wordCounts, int64(i+10000))
		v := make([]float32, 0, 2)
		for j := 0; j < 2; j++ {
			v = append(v, rand.Float32())
		}
		bookIntros = append(bookIntros, v)
	}
	idColumn := entity.NewColumnInt64("book_id", bookIDs)
	wordColumn := entity.NewColumnInt64("word_count", wordCounts)
	introColumn := entity.NewColumnFloatVector("book_intro", 2, bookIntros)
	// insert
	_, err := milvusClient.Insert(
		ctx,         // ctx
		collection,  // CollectionName
		"",          // partitionName
		idColumn,    // columnarData
		wordColumn,  // columnarData
		introColumn, // columnarData
	)
	if err != nil {
		log.Fatal("failed to insert data:", err.Error())
	}
	fmt.Printf("Successfully inserted data into %s\n", collection)
}

func CreateIndex(milvusClient client.Client, collection string, fieldName string, level entity.MetricType, nlist int, ctx context.Context) error {
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		level, // metricType
		nlist, // ConstructParams
	)
	if err != nil {
		log.Fatal("fail to create ivf flat index parameter:", err.Error())
	}
	err = milvusClient.CreateIndex(
		context.Background(), // ctx
		collection,           // CollectionName
		fieldName,            // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		log.Fatal("fail to create index:", err.Error())
	}
	return nil
}

func LoadCollection(milvusClient client.Client, collection string, ctx context.Context) error {
	err := milvusClient.LoadCollection(
		context.Background(), // ctx
		collection,           // CollectionName
		false,                // async
	)
	if err != nil {
		log.Fatal("fail to load collection:", err.Error())
	}
	fmt.Println("Successfully loaded collection:", collection)
	return nil
}

func ConductSearch(milvusClient client.Client, collection string, outputFields []string, queryVectors []float32, topK int, ctx context.Context) error {
	sp, _ := entity.NewIndexFlatSearchParam()

	searchResult, err := milvusClient.Search(
		context.Background(), // ctx
		collection,           // CollectionName
		[]string{},           // partitionNames
		"",                   // expr
		outputFields,         // outputFields
		[]entity.Vector{entity.FloatVector(queryVectors)}, // vectors
		"book_intro", // vectorField
		entity.L2,    // metricType
		topK,         // topK
		sp,           // sp
	)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}

	fmt.Printf("%#v\n", searchResult)
	for _, sr := range searchResult {
		fmt.Println(sr.IDs)
		fmt.Println(sr.Scores)
		fmt.Printf("%#v\n", sr.Fields)
	}
	return nil
}
