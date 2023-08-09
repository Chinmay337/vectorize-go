package vectordb

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type SearchResult struct {
	IDs    *entity.ColumnVarChar
	Scores []float32
	Fields []*entity.ColumnVarChar
}

func SearchIndexFromCollection(milvusClient client.Client, collection string, queryField string, queryVectors []entity.Vector, outputFields []string, topK int, ctx context.Context) error {
	sp, _ := entity.NewIndexFlatSearchParam()

	searchResult, err := milvusClient.Search(
		context.Background(), // ctx
		collection,           // CollectionName
		[]string{},           // partitionNames
		"",                   // expr
		outputFields,         // outputFields
		queryVectors,         // vectors
		queryField,           // vectorField
		entity.L2,            // metricType
		topK,                 // topK
		sp,                   // sp
	)
	if err != nil {
		log.Fatal("fail to search collection:", err.Error())
	}

	fmt.Printf("Raw Search Result\n%#v\n\n", searchResult)

	for _, sr := range searchResult {
		fmt.Println("IDs ", sr.IDs)

		numResults := sr.ResultCount
		fmt.Println("Number of Results: ", numResults)
		for i := 0; i < numResults; i++ {
			ids, _ := sr.IDs.GetAsString(i)
			fmt.Printf("Result %d: %s Score: %f\n", i, ids, sr.Scores[i])
		}

		// fmt.Println("Scores ", sr.Scores)
		fmt.Printf("Fields:\n%#v\n", sr.Fields)

		for _, field := range sr.Fields {
			fmt.Println(field.FieldData().GetScalars().Data)
		}
	}

	return nil
}
