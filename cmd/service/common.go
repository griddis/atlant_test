package service

import "time"

type FetchRequest struct {
	Files []string `json:"files"`
}

type FetchResponse struct {
	Status string `json:"status"`
}

type Fetch struct {
}

type ListRequest struct {
	Limiter *Limiter         `json:"limiter,omitempty"`
	Sorter  map[string]int32 `json:"sorter,omitempty"`
}

type ListResponse struct {
	Status    string       `json:"status"`
	ListItems []*ListItems `json:"listItems"`
}

type ListItems struct {
	Id      string
	Name    string
	Price   float32
	Counter uint32
	Date    time.Time
}
type Limiter struct {
	Offsetbyid string `json:"offsetById,omitempty"`
	Limit      int64  `json:"limit,omitempty"`
}
