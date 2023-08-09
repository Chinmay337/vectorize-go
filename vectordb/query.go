package vectordb

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

func QueryCollection(milvusClient client.Client, collection string, expr string, outputFields []string, ctx context.Context) {
	queryResult, err := milvusClient.Query(
		context.Background(), // ctx
		collection,           // CollectionName
		[]string{},           // PartitionName
		expr,                 // expr
		outputFields,         // OutputFields
	)
	if err != nil {
		log.Fatal("fail to query collection:", err.Error())
	}

	fmt.Printf("Raw Query Result\n%#v\n\n", queryResult)

	for _, qr := range queryResult {

		fmt.Println("Column Name:", qr.Name())

		numResults := qr.Len()
		fmt.Printf("Query Returned %d results\n", numResults)

		for i := 0; i < numResults; i++ {
			val, _ := qr.GetAsString(i)
			fmt.Printf("Result %d: %s\n", i, val)
		}

		// fmt.Println("Field Data : ", qr.FieldData())
		// fields := qr.FieldData().GetScalars().Data

	}
}
