# Vector DB

## Running

```bash
# Start Milvus DB and MinIO
docker-compose up -d

# Run the server
go run .

```

## Tests and Benchmarks

```bash
# Run all tests in package
go test -v ./...

# Run all benchmarks in package
go test -bench ./...

# View all benchmarks in project
go test -list '.*' ./...

go test -bench=.

# Memory and CPU
go test -bench=. -benchmem -cpuprofile cpu.out -memprofile mem.out

# Use pprof to view the metrics
go tool pprof -http=:8080 cpu.out
go tool pprof -http=:8080 mem.out

```
