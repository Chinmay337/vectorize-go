package vectorize

import (
	"log"
	"os"

	"github.com/ynqa/wego/pkg/model/modelutil/vector"
	"github.com/ynqa/wego/pkg/model/word2vec"

	"github.com/ynqa/wego/pkg/embedding"
	"github.com/ynqa/wego/pkg/search"
)

func Train() error {
	model, err := word2vec.New(
		word2vec.Window(5),
		word2vec.Model(word2vec.Cbow),
		word2vec.Optimizer(word2vec.NegativeSampling),
		word2vec.NegativeSampleSize(5),
		word2vec.Verbose(),
	)
	if err != nil {
		// failed to create word2vec.
	}

	input, _ := os.Open("string-vectors/input")
	defer input.Close()
	if err = model.Train(input); err != nil {
		// failed to train.
	}

	// write word vector to a file
	output, err := os.Create("string-vectors/word_vector.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	// Save Trained Model to Disk
	model.Save(output, vector.Agg)

	// write word vector to stdout
	// model.Save(os.Stdin, vector.Agg)

	// print vectors for all words
	// for _, word := range []string{"english", "hindi", "dog", "cat", "sanskrit", "memes", "lithuania", "darwin", "marathi", "german", "italian"} {
	// 	vector := model.WordVector(word)
	// 	if err != nil {
	// 		fmt.Printf("Error getting vector for word %s: %s\n", word, err)
	// 		continue
	// 	}
	// 	fmt.Printf("Vector for word %s: %v\n", word, vector)
	// }

	return nil
}

func QueryVector(word string) error {
	input, err := os.Open("string-vectors/word_vector.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer input.Close()
	embs, err := embedding.Load(input)
	if err != nil {
		log.Fatal(err)
	}
	searcher, err := search.New(embs...)
	if err != nil {
		log.Fatal(err)
	}
	neighbors, err := searcher.SearchInternal(word, 10)
	if err != nil {
		log.Fatal(err)
	}
	neighbors.Describe()
	return nil
}

// func VectorizeString(word string) ([]float64, error) {
// 	model, err := word2vec.Load("path/to/trained/model")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	vector, err := model.WordVector(word)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(vector) // print the vector to the console

// 	return vector, nil
// }
