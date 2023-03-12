package storage

import (
	"fmt"
	"os"
	"path"
)

type FileStorage struct {
	root string
}

func NewFileStorage(root string) (*FileStorage, error) {
	if _, err := os.Stat(root); err != nil {
		err := os.MkdirAll(root, 0755)
		if err != nil {
			return nil, fmt.Errorf("Error creating root directory: %s", err)
		}
	}

	return &FileStorage{
		root: root,
	}, nil
}

func (s *FileStorage) Get(filePath string) (string, error) {
	fullPath := path.Join(s.root, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content := make([]byte, 10000)
	_, err = file.Read(content)
	if err != nil {
		return "", err
	}

	return string(content), nil

}

func (s *FileStorage) Set(filePath string, value string) error {
	fullPath := path.Join(s.root, filePath)
	err := os.MkdirAll(path.Dir(fullPath), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(value)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) Delete(filePath string) error {
	return nil
}
