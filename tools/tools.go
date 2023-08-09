package tools

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var pFlag = flag.Bool("perf", false, "Start the performance server")

func EnablePerformanceServerIfFlag() {
	flag.Parse()

	if *pFlag {
		go func() {
			log.Println("Starting performance server on localhost:6060")
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

func ConnectVectorDB() client.Client {
	log.Print("Connecting to VectorDB...")
	var err error
	done := make(chan bool)

	go func() {
		time.Sleep(10 * time.Second)
		select {
		case <-done:
			return
		default:
			// Panic if connection unsuccessful after 5s
			log.Print("failed to connect to Milvus. Make sure to run\ndocker-compose up\n")
			os.Exit(1)
		}
	}()

	milvusClient, err := client.NewGrpcClient( // Max 65,536 connections
		context.Background(), // ctx
		"localhost:19530",    // addr
	)
	close(done)
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
