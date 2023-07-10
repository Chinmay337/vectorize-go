# Milvus Intro 


- Milvus 

- https://milvus.io/docs/manage_connection.md


## Key Concepts 

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
docker run --name milvusui -d -p 8000:3000 -e HOST_URL=http://172.27.169.49:8000 -e MILVUS_URL=172.27.169.49:19530 zilliz/attu:latest

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


