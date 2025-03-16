package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/sur1k1/go-url-shortener/internal/models"
	"github.com/sur1k1/go-url-shortener/internal/repository"
	"go.uber.org/zap"
)

type MemStorage struct {
	URLs map[string]models.URLData
	c    sync.RWMutex
	file *os.File
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewStorage(log *zap.Logger, path string) (*MemStorage, error) {
	const op = "storage.NewStorage"

	file, err := loadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to laod file: %v", op, err)
	}

	s := &MemStorage{
		URLs: map[string]models.URLData{},
		c: sync.RWMutex{},
		file: file,
	}

	fileStat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get file stat: %v", op, err)
	}

	if fileStat.Size() > 0 {
		if err := s.restoreStorage(); err != nil {
			return nil, fmt.Errorf("%s: cannot restore data urls: %v", op, err)
		}
	}

	return s, nil
}

func loadFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE | os.O_RDWR | os.O_APPEND, 0644)
}

func (s *MemStorage) restoreStorage() error {
	const op = "storage.restoreStorage"

	scanner := bufio.NewScanner(s.file)

	var urlData models.URLData

	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &urlData)
		if err != nil {
			return fmt.Errorf("%s: failed to unmarshal event: %v", op, err)
		}

		s.URLs[urlData.ShortURL] = urlData
	}

	return nil
}

func (s *MemStorage) GetURL(shortURL string) (*models.URLData, error) {
	const op = "storage.GetURL"

	s.c.RLock()
	defer s.c.RUnlock()

	urlData, ok := s.URLs[shortURL]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrURLNotFound)
	}

	return &urlData, nil
}

func (s *MemStorage) SaveURL(urlData *models.URLData) error {
	s.c.Lock()
	defer s.c.Unlock()

	urlData.UUID = strconv.Itoa(len(s.URLs) + 1)
	s.URLs[urlData.ShortURL] = *urlData
	s.writeFile(*urlData)

	return nil
}

func (s *MemStorage) writeFile(urlData models.URLData) error {
	const op = "storage.writeFile"

	jsonData, err := json.Marshal(urlData)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal data: %v", op, err)
	}

	jsonData = append(jsonData, '\n')

	_, err = s.file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("%s: failed to write file: %v", op, err)
	}

	return nil
}

func (s *MemStorage) Close() error {
	return s.file.Close()
}