package service

import (
	"os"
	"strings"
)

type FileManager interface {
	Create(URI string) (*os.File, error)
	Delete(URI string) error
}

type FileManagerImpl struct{}

func NewFileManager() FileManager {
	return &FileManagerImpl{}
}

func (f *FileManagerImpl) Create(URI string) (*os.File, error) {
	filename := nameFromUri(URI)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *FileManagerImpl) Delete(URI string) error {
	filename := nameFromUri(URI)
	return os.Remove(filename)
}

func nameFromUri(URI string) string {
	URIParts := strings.Split(URI, "/")
	length := len(URIParts)
	filename := URIParts[length-1]
	return filename
}
