package repository

import (
	"context"
	"errors"
	"time"
)

type Repository interface {
	//Connect(ctx context.Context, database, coll string) error
	UpdatePrice(ctx context.Context, obj ProductPrice) error
	ListPrice(ctx context.Context, sorter map[string]int32, limiter Limiter) ([]*ProductPrice, error)
	Close(ctx context.Context) error
}

var Noresult = errors.New("no result")

type ProductPrice struct {
	ID      string    `json:"id" bson:"_id,omitempty"`
	Name    string    `csv:"name" json:"name,omitempty" bson:"name,omitempty"`
	Price   float32   `csv:"price" json:"price,omitempty" bson:"price,omitempty"`
	Date    time.Time `csv:"date" json:"date,omitempty" bson:"date,omitempty"`
	Counter uint32    `csv:"counter" json:"counter,omitempty" bson:"counter,omitempty"`
}

type Limiter struct {
	Limit      int64
	Offsetbyid string
}
