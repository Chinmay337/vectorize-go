package vectorize

import (
	"testing"

	"milvus/errors"
)

const (
	validInputPath = "../tests/testdata"
	emptyInputPath = "../tests/empty"
	validModelPath = "../tests/word_vector.txt"
	unknownPath    = "imaginary/path/to/invalid/input.txt"
)

func TestTrain(t *testing.T) {
	tests := []struct {
		name          string
		inputPath     string
		outputPath    string
		expectedErrFn func(error) bool // This function should return true if the error matches expectations.
	}{
		{
			name:       "Valid Input",
			inputPath:  validInputPath,
			outputPath: validModelPath,
			expectedErrFn: func(err error) bool {
				return err == nil
			},
		},
		{
			name:       "Invalid Input - Passing unknown file path",
			inputPath:  unknownPath,
			outputPath: validModelPath,
			expectedErrFn: func(err error) bool {
				return errors.IsFileError(err, "FileNotFound")
			},
		},
		{
			name:       "Invalid Input - Passing empty file",
			inputPath:  emptyInputPath,
			outputPath: validModelPath,
			expectedErrFn: func(err error) bool {
				return errors.IsFileError(err, "FileEmpty")
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Train(tt.inputPath, tt.outputPath)
			if !tt.expectedErrFn(err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestQueryVector(t *testing.T) {
	tests := []struct {
		name          string
		word          string
		inputPath     string
		expectedErrFn func(error) bool
	}{
		{
			name:      "Valid Query",
			word:      "cat",
			inputPath: validModelPath,
			expectedErrFn: func(err error) bool {
				return err == nil
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := QueryVector(tt.word, tt.inputPath)
			if !tt.expectedErrFn(err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
