package vectordb

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type CollectionBuilder struct {
	name          string
	description   string
	fields        []*entity.Field
	enableDynamic bool
	shardNum      int32
}

func NewCollectionBuilder() *CollectionBuilder {
	return &CollectionBuilder{
		fields:   make([]*entity.Field, 0),
		shardNum: 2, // default value
	}
}

func (cb *CollectionBuilder) WithName(name string) *CollectionBuilder {
	cb.name = name
	return cb
}

func (cb *CollectionBuilder) WithDescription(desc string) *CollectionBuilder {
	cb.description = desc
	return cb
}

func (cb *CollectionBuilder) WithFields(fields ...*entity.Field) *CollectionBuilder {
	cb.fields = append(cb.fields, fields...)
	return cb
}

// Helper function for creating VarChar fields
func NewFieldVarChar(name string, maxLength int, primaryKey bool, autoID bool) *entity.Field {
	return &entity.Field{
		Name:     name,
		DataType: entity.FieldTypeVarChar,
		TypeParams: map[string]string{
			"max_length": fmt.Sprintf("%d", maxLength),
		},
		PrimaryKey: primaryKey,
		AutoID:     autoID,
	}
}

// Helper function for creating FloatVector fields
func NewFieldFloatVector(name string, dim int) *entity.Field {
	return &entity.Field{
		Name:     name,
		DataType: entity.FieldTypeFloatVector,
		TypeParams: map[string]string{
			"dim": fmt.Sprintf("%d", dim),
		},
	}
}

func (cb *CollectionBuilder) Create(milvusClient client.Client, ctx context.Context) error {
	schema := &entity.Schema{
		CollectionName:     cb.name,
		Description:        cb.description,
		Fields:             cb.fields,
		EnableDynamicField: cb.enableDynamic,
	}

	collections, err := milvusClient.ListCollections(ctx)
	if err != nil {
		log.Fatal("failed to list collections:", err.Error())
	}
	for _, collection := range collections {
		if collection.Name == cb.name {
			fmt.Printf("Collection %s already exists\n", cb.name)
			return nil
		}
	}
	err = milvusClient.CreateCollection(ctx, schema, cb.shardNum)
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}
	fmt.Printf("Successfully created collection %s\n", cb.name)
	return nil
}
