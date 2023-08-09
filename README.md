# Vector DB - Concepts and Use Cases

- Model that supports unknown words https://fasttext.cc/docs/en/faqs.html

- Can use `Whisper.cpp` to add Subtitles to Videos - it can also run in `WASM` in the browser

## Key Concepts

- Basics

</br>

```
Milvus is an open-source vector database that provides state-of-the-art similarity search and analysis capabilities.

1. Create Collection
2. Insert Vectors
3. Create an Index on the Collection
4. Load Collection
5. Perform Search on Collection

```

</br>

- Collections

</br>

```
A collection is a set of vectors that share the same schema.
Collection must have a Primary Key
The PK can either be VARCHAR or INT64

We need a model to generate Vectors

$ go get -u github.com/ynqa/wego

$ go install github.com/ynqa/wego

$ wego -h

```

</br>

- Indexes

</br>

```go
Metric Type (L2) = Type of metrics used to measure the similarity of vectors.

Metric Type can be :
For floating point vectors:
L2 (Euclidean distance)
IP (Inner product)
For binary vectors:
JACCARD (Jaccard distance)
TANIMOTO (Tanimoto distance)
HAMMING (Hamming distance)
SUPERSTRUCTURE (Superstructure)
SUBSTRUCTURE (Substructure)

Construct Params - In Memory Index / On Disk index Size
1024 is a default

idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024, // ConstructParams
	)

idx : Pass above index to CreateIndex function
fieldName : Field from collection to Index for Searches

err = milvusClient.CreateIndex(
		context.Background(),         // ctx
		"collectionname",             // CollectionName
		"fieldFromCollectionToIndex", // fieldName
		idx,                          // entity.Index
		false,                        // async
	)

L2 (Euclidean distance)
IP (Inner product)
```

</br>

- Querying - query items/fields that match some condition

</br>

```go
// We need to provide a query Expression
vectordb.QueryCollection(client, "words", "word not in ['cat', 'dog']", []string{"word"}, ctx)

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

	fmt.Printf("%#v\n", queryResult)
	for _, qr := range queryResult {
		fmt.Printf("%#v\n", qr)
	}
}

```

</br>

- Searching - search for similarity

</br>

```go
// Load Collection
err := milvusClient.LoadCollection(
  context.Background(),   // ctx
  "book",                 // CollectionName
  false                   // async
)
if err != nil {
  log.Fatal("failed to load collection:", err.Error())
}

// Prepare search parameters

sp, _ := entity.NewIndexFlatSearchParam( // NewIndex*SearchParam func
	10,                                  // searchParam
)

// Conduct a Vector Similarity Search

searchResult, err := milvusClient.Search(
	context.Background(),                    // ctx
	"book",                                  // CollectionName
	[]string{},                              // partitionNames
	"",                                      // expr
	[]string{"book_id"},                     // outputFields
	[]entity.Vector{entity.FloatVector([]float32{0.1, 0.2})}, // vectors
	"book_intro",                            // vectorField
	entity.L2,                               // metricType
	2,                                       // topK
	sp,                                      // sp
)
if err != nil {
	log.Fatal("fail to search collection:", err.Error())
}

// Accessing Search Results

	for _, sr := range searchResult {
		fmt.Println("IDs ", sr.IDs)

		numResults := sr.ResultCount
		fmt.Println("Number of Results: ", numResults)
		for i := 0; i < numResults; i++ {
			ids, _ := sr.IDs.GetAsString(i)
			fmt.Printf("Result %d: %s Score: %f\n", i, ids, sr.Scores[i])
		}


		fmt.Printf("Fields:\n%#v\n", sr.Fields)

		for _, field := range sr.Fields {
			fmt.Println(field.FieldData().GetScalars().Data)
		}
	}


IDs  &{{}  [dog cat word2]}

Number of Results:  3
Result 0: dog Score: 0.000000
Result 1: cat Score: 0.010000
Result 2: word2 Score: 0.170000
Fields:
client.ResultSet{(*entity.ColumnVarChar)(0xc000609b90)}
&{data:"dog" data:"cat" data:"word2" }

Lower Score means better match - 0.0 is a perfect match


```

</br>

```bash
# Download Milvus
wget https://github.com/milvus-io/milvus/releases/download/v2.2.11/milvus-standalone-docker-compose.yml -O docker-compose.yml

# Spin it up
sudo docker-compose up -d

# Check if Containers are up
sudo docker-compose ps

# Connect to Milvus
docker port milvus-standalone 19530/tcp

```

- UI for Milvus - `Attu`

```bash
# UI for Milvus
docker run -p 8000:3000 -e HOST_URL=http://{ your machine IP }:8000 -e MILVUS_URL={your machine IP}:19530 zilliz/attu:latest

# Getting IP info
ip addr show eth0 | grep inet | awk '{ print $2; }' | sed 's/\/.*$//'

# Command
docker run --name milvusui -d -p 8000:3000 -e HOST_URL=http://machip:8000 -e MILVUS_URL=machip:19530 zilliz/attu:latest

# Visit localhost:8000 to view UI

```

## Go

```bash
# Installation WSL2
wget https://dl.google.com/go/go1.20.5.linux-amd64.tar.gz
sudo tar -xvf go1.20.5.linux-amd64.tar.gz
sudo mv go /usr/local/go-1.20.5
source ~/.bashrc
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Usage
golangci-lint run # Run linter

```

```bash
Awesome Tools
- benchcmp  # Compare benchmark results
- benchstat # Analyze benchmark results
- go tool pprof # Analyze CPU and memory profile

- go-callvis # Visualize call graph
- go-torch # Visualize CPU profile
- go tool trace # Analyze trace
- go tool trace -http=:8080 trace.out # Analyze trace
```

- Go Milvus

</br>

```go
go get -u github.com/milvus-io/milvus-sdk-go/v2
go get github.com/milvus-io/milvus-sdk-go/v2/client@v2.2.6
```

</br>

- Go Profiling

```bash
#Profiling

go install github.com/uber/go-torch

go-torch -u http://localhost:6060/

# Needs this at right path
sudo git clone https://github.com/brendangregg/FlameGraph.git

# To find out where to clone so it is always available without adding to path
which go # /usr/local/go-1.20.5/bin/go

# So we clone Flamegraph to /usr/local/
cd /usr/local/ && sudo git clone https://github.com/brendangregg/FlameGraph.git


```
