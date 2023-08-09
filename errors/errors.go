package errors

import "fmt"

type FileError struct {
	Type    string
	Path    string
	Message string
}

func (fe *FileError) Error() string {
	return fmt.Sprintf("%s: %s - %s", fe.Type, fe.Path, fe.Message)
}

func IsFileError(err error, action string) bool {
	fe, ok := err.(*FileError)
	if !ok {
		return false
	}
	return fe.Type == action
}

var (
	FileNotFound = func(path string, err error) error {
		return &FileError{Type: "FileNotFound", Path: path, Message: err.Error()}
	}
	FileEmpty = func(path string, err error) error {
		return &FileError{Type: "FileEmpty", Path: path, Message: err.Error()}
	}
	FileCreationErr = func(path string, err error) error {
		return &FileError{Type: "FileCreationError", Path: path, Message: err.Error()}
	}
	FileLoadingError = func(path string, err error) error {
		return &FileError{Type: "FileLoadingError", Path: path, Message: err.Error()}
	}

	ModelSearchError = func(path string, err error) error {
		return &FileError{Type: "FileLoadingError", Path: path, Message: err.Error()}
	}

	ModelLoadingError = func(path string, err error) error {
		return &FileError{Type: "FileLoadingError", Path: path, Message: err.Error()}
	}
)
