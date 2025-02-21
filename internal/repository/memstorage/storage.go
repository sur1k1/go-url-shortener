package storage

import "sync"

type MemStorage struct {
	URLs map[string]string
	c    sync.RWMutex
}

func NewStorage() *MemStorage {
	return &MemStorage{
		URLs: map[string]string{},
		c: sync.RWMutex{},
	}
}

func (s *MemStorage) GetURL(shortURL string) (string, bool) {
	s.c.RLock()
	defer s.c.RUnlock()

	url, ok := s.URLs[shortURL]
	
	return url, ok
}

func (s *MemStorage) SaveURL(shortURL string, originalURL string) {
	s.c.Lock()
	defer s.c.Unlock()

	s.URLs[shortURL] = originalURL
}