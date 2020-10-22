package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/griddis/atlant_test/tools/logging"
	"github.com/imdario/mergo"
)

type inMemoryStore struct {
	logger *logging.Logger
	mx     sync.RWMutex
	items  prices
}

var NoResults = errors.New("no results")

type prices map[string]*ProductPrice

func NewInMemoryStore(ctx context.Context) (Repository, error) {
	logger := logging.FromContext(ctx)
	return &inMemoryStore{
		logger: logger,
		items:  make(prices, 0),
	}, nil
}

func (s *inMemoryStore) UpdatePrice(ctx context.Context, obj ProductPrice) error {

	s.logger.Debug("msg", "UpdatePrice")
	cur := ProductPrice{}
	hash := md5.Sum([]byte(obj.Name))
	index := hex.EncodeToString(hash[:])
	cur, err := s.GetByID(ctx, index)

	if err == nil {
		if err := mergo.Merge(&cur, obj, mergo.WithOverride); err != nil {
			return errors.New("merge obj error")
		}
		cur.Counter++
	} else {
		cur = obj
		cur.Counter = 0
	}

	if cur.ID == "" {
		cur.ID = uuid.New().String()
	}

	cur.Date = time.Now()

	s.mx.Lock()
	s.items[index] = &cur
	s.mx.Unlock()
	return nil
}
func (s *inMemoryStore) ListPrice(ctx context.Context, sorter map[string]int32, limiter Limiter) ([]*ProductPrice, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return nil, nil
}

func (s *inMemoryStore) GetByID(ctx context.Context, index string) (ProductPrice, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if c, ok := s.items[index]; ok {
		return *c, nil
	}

	return ProductPrice{}, NoResults
}

func (s *inMemoryStore) Close(ctx context.Context) error {
	return nil
}
