```

 ██████╗  ██████╗     ███╗   ███╗██╗██╗    ██╗   ██╗██╗   ██╗███████╗
██╔════╝ ██╔═══██╗    ████╗ ████║██║██║    ██║   ██║██║   ██║██╔════╝
██║  ███╗██║   ██║    ██╔████╔██║██║██║    ██║   ██║██║   ██║███████╗
██║   ██║██║   ██║    ██║╚██╔╝██║██║██║    ╚██╗ ██╔╝██║   ██║╚════██║
╚██████╔╝╚██████╔╝    ██║ ╚═╝ ██║██║███████╗╚████╔╝ ╚██████╔╝███████║
 ╚═════╝  ╚═════╝     ╚═╝     ╚═╝╚═╝╚══════╝ ╚═══╝   ╚═════╝ ╚══════╝

```

# Vector DB Interface

- Running

```bash
# Start Milvus DB and MinIO
docker-compose up -d

# Run the server
go run .

# Run the server with performance profiling enabled
go run . -perf

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

- Model that supports unknown words https://fasttext.cc/docs/en/faqs.html

- Can use `Whisper.cpp` to add Subtitles to Videos - it can also run in `WASM` in the browser
