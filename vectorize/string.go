package vectorize

import (
	"fmt"
	"os"

	stdErrors "errors"
	"milvus/errors"

	"github.com/ynqa/wego/pkg/model/modelutil/vector"
	"github.com/ynqa/wego/pkg/model/word2vec"

	"github.com/ynqa/wego/pkg/embedding"
	"github.com/ynqa/wego/pkg/search"
)

func Train(inputPath string, outputPath string) error {
	fmt.Printf("Training data from : %s\n", inputPath)

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		if stdErrors.Is(err, os.ErrNotExist) {
			return errors.FileNotFound(inputPath, err) // here 'errors' refers to your custom package
		}
		return errors.FileLoadingError(inputPath, err)
	}

	// Check if the input file is empty.
	if fileInfo.Size() == 0 {
		return errors.FileEmpty(inputPath, stdErrors.New("Empty File"))
	}

	model, err := word2vec.New(
		word2vec.Window(5),
		word2vec.Model(word2vec.Cbow),
		word2vec.Optimizer(word2vec.NegativeSampling),
		word2vec.NegativeSampleSize(5),
	)
	if err != nil {
		fmt.Printf("Failed to Load model: %s\n", err)
		return errors.ModelLoadingError(inputPath, err)

	}

	input, _ := os.Open(inputPath)
	if err != nil {
		return errors.FileLoadingError(inputPath, err)
	}
	defer input.Close()
	if err = model.Train(input); err != nil {
		// failed to train.
		fmt.Printf("Error training model: %s\n", err)
	}

	// write word vector to a file
	output, err := os.Create(outputPath)
	if err != nil {
		return errors.FileCreationErr(inputPath, err)
	}
	defer output.Close()

	// Save Trained Model to Disk
	model.Save(output, vector.Agg)

	fmt.Printf("Successfully trained and saved model to %s\n", outputPath)

	return nil
}

func QueryVector(word string, inputPath string) error {
	fmt.Printf("Querying similarity for word: %s\n", word)

	input, err := os.Open(inputPath)
	if err != nil {
		return errors.FileCreationErr(inputPath, err)
	}
	defer input.Close()
	embs, err := embedding.Load(input)
	if err != nil {
		return errors.FileLoadingError(inputPath, err)
	}
	searcher, err := search.New(embs...)
	if err != nil {
		return errors.ModelSearchError(inputPath, err)
	}
	neighbors, err := searcher.SearchInternal(word, 10)
	if err != nil {
		return errors.ModelSearchError(inputPath, err)
	}
	neighbors.Describe()
	return nil
}
