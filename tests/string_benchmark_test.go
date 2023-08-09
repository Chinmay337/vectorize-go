package tests

import (
	"testing"

	"milvus/vectorize"
)

func BenchmarkTrain(b *testing.B) {
	inputPath := "mockdata/string-vectors/input.txt"
	outputPath := "mockdata/string-vectors/word_vector.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vectorize.Train(inputPath, outputPath)
	}
}

func BenchmarkQueryVector(b *testing.B) {
	word := "benchmarkWord"
	inputPath := "string-vectors/input.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vectorize.QueryVector(word, inputPath)
	}
}
