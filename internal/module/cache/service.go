package cache

import (
	"fmt"
	"time"
)

type ICacheRepository interface {
	SetValue(string, interface{}, time.Duration) error
	GetValue(string) (string, error)
}

type CacheService struct {
	CacheRepo ICacheRepository
}

func NewCacheService(cacheRepo ICacheRepository) *CacheService {
	return &CacheService{
		CacheRepo: cacheRepo,
	}
}

func (s *CacheService) SetValue(key string, value interface{}, expiration time.Duration) error {
	err := s.CacheRepo.SetValue(key, value, expiration)
	if err != nil {
		return fmt.Errorf("failed to set value: %s", err.Error())
	}
	return nil
}

func (s *CacheService) GetValue(key string) (string, error) {
	return s.CacheRepo.GetValue(key)
}
