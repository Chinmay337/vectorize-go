package vectorize

import (
	"testing"
)

func BenchmarkTrain(b *testing.B) {
	inputPath := "string-vectors/input.txt"
	outputPath := "string-vectors/word_vector.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Train(inputPath, outputPath)
	}
}

func BenchmarkQueryVector(b *testing.B) {
	word := "benchmarkWord"
	inputPath := "string-vectors/input.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QueryVector(word, inputPath)
	}
}
