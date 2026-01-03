package service

import (
	"context"
	"errors"
	"github.com/lhducc/bookmark-management/internal/repository"
	"github.com/lhducc/bookmark-management/pkg/stringutils"
)

const (
	urlCodeLength = 7
	maxRetry      = 5
)

//go:generate mockery --name ShortenUrl --filename urlstorage.go
type ShortenUrl interface {
	ShortenUrl(ctx context.Context, url string, exp int) (string, error)
}

type shortenUrl struct {
	repo repository.UrlStorage
}

func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrl{repo: repo}
}

// ShortenUrl shortens a given URL and returns a shortened URL code.
// The method generates a random URL code of length urlCodeLength, stores the given URL with the generated URL code in the repository, and returns the generated URL code.
// If an error occurs while generating the URL code, it returns an empty string and the error immediately.
// If an error occurs while storing the URL in the repository, it returns an empty string and the error immediately.
// The returned URL code is a string of length urlCodeLength, and does not contain any whitespace or special characters.
// The URL code is case-sensitive and can be used to retrieve the original URL from the repository.
func (s *shortenUrl) ShortenUrl(ctx context.Context, url string, exp int) (string, error) {
	for i := 0; i < maxRetry; i++ {
		urlCode, err := stringutils.GenerateCode(urlCodeLength)
		if err != nil {
			return "", err
		}

		ok, err := s.repo.StoreURLIfNotExists(ctx, urlCode, url, exp)
		if err != nil {
			return "", err
		}

		if ok {
			return urlCode, nil
		}
	}
	return "", errors.New("failed to shorten URL")
}
