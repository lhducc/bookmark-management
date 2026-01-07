package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	urlExpTime = 24 * time.Hour
)

//go:generate mockery --name=UrlStorage --filename urlstorage.go
type UrlStorage interface {
	StoreURL(ctx context.Context, code, url string) error
	GetURL(ctx context.Context, code string) (string, error)
	StoreURLIfNotExists(ctx context.Context, code, url string, exp int) (bool, error)
}
type urlStorage struct {
	c *redis.Client
}

func NewUrlStorage(c *redis.Client) UrlStorage {
	return &urlStorage{c: c}
}

// StoreURL stores a URL in the repository with a given code and expiration time.
// The method takes a context, a code, and a URL as input parameters.
// It stores the URL in the repository with the given code and expiration time, and returns an error if there is an issue storing the URL.
func (s *urlStorage) StoreURL(ctx context.Context, code, url string) error {
	return s.c.Set(ctx, code, url, urlExpTime).Err()
}

// GetURL retrieves a URL from the repository using a given code.
// The method takes a context and a code as input parameters.
// It returns the URL associated with the given code, and an error if there is an issue retrieving the URL.
func (s *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return s.c.Get(ctx, code).Result()
}

func (s *urlStorage) StoreURLIfNotExists(ctx context.Context, code, url string, exp int) (bool, error) {
	expDuration := urlExpTime
	if exp > 0 {
		expDuration = time.Duration(exp) * time.Second
	}

	ok, err := s.c.SetNX(ctx, code, url, expDuration).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}
