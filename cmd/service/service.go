package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/griddis/atlant_test/configs"
	"github.com/griddis/atlant_test/pkg/crawler"
	"github.com/griddis/atlant_test/pkg/repository"
	"github.com/griddis/atlant_test/tools/logging"
	"github.com/pkg/errors"
)

type service struct {
	config     *configs.Config
	repository repository.Repository
	client     *http.Client
}

func NewService(ctx context.Context, cfg *configs.Config, repo repository.Repository, client *http.Client) Service {
	return &service{
		config:     cfg,
		repository: repo,
		client:     client,
	}
}

func (s *service) Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With("svc", "main", "func", "Fetch")

	countFiles := len(req.Files)
	concurrency := s.config.Crawler.Concurrency
	if countFiles < s.config.Crawler.Concurrency {
		concurrency = countFiles
	}
	crawlerService := crawler.NewCrawler(ctx, concurrency, s.client)

	go func() {
		select {
		case <-ctx.Done():
			logger.Debug("msg", "request done")
			crawlerService.Close()
		}
	}()

	go func() {
		for i := 0; i < countFiles; i++ {
			crawlerService.Input <- req.Files[i]
		}
	}()

	for i := 0; i < countFiles; i++ {
		str := <-crawlerService.Output
		if len(str) != 0 {
			logger.Debug("status", "success", "req", str)

			r := csv.NewReader(strings.NewReader(str))

			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, errors.Wrap(err, "read error csv")
				}

				if len(record[0]) == 0 {
					logger.Error("name", "is empty")
					continue
				}

				price, err := strconv.ParseFloat(record[1], 64)
				if err != nil {
					logger.Error("price", "parse error")
					continue
				}

				if price <= 0 {
					logger.Error("price", "only more than 0")
					continue
				}

				product := repository.ProductPrice{Name: record[0], Price: float32(price)}

				err = s.repository.UpdatePrice(ctx, product)
				if err != nil {
					return nil, errors.Wrap(err, fmt.Sprintf("err update %+v", product))
				}
			}
		}
	}

	return &FetchResponse{Status: "OK"}, nil
}

func (s *service) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With("method", "list")
	limiter := repository.Limiter{
		Offsetbyid: req.Limiter.Offsetbyid,
		Limit:      req.Limiter.Limit,
	}
	logger.Debug("sorter", fmt.Sprintf("%+v", req.Sorter))
	products, err := s.repository.ListPrice(ctx, req.Sorter, limiter)
	if err != nil {
		return &ListResponse{Status: "error"}, errors.Wrap(err, "dont listing products")
	}
	var list []*ListItems
	for _, p := range products {
		item := ListItems{
			Id:      p.ID,
			Name:    p.Name,
			Price:   p.Price,
			Date:    p.Date,
			Counter: p.Counter,
		}
		list = append(list, &item)
	}
	return &ListResponse{Status: "OK", ListItems: list}, nil
}

func (s *service) _Fetch(ctx context.Context, req *FetchRequest) (*FetchResponse, error) {
	logger := logging.FromContext(ctx)
	logger = logger.With("svc", "main", "func", "Fetch")

	countFiles := len(req.Files)
	if countFiles <= 0 {
		return nil, errors.New("No fetch files")
	}

	for _, v := range req.Files {
		req, _ := http.NewRequest("GET", v, nil)
		req = req.WithContext(ctx)
		tr := &http.Transport{}
		client := &http.Client{}
		client.Transport = tr
		errChan := make(chan error)

		go func() {
			resp, err := client.Do(req)
			if err != nil {
				errChan <- errors.Wrap(err, "Request error")
				return
			}

			if resp.StatusCode != http.StatusOK {
				logger.Error("err", "status code is not 200")
				errChan <- errors.New("status code is not 200")
				return
			}

			defer resp.Body.Close()
			respData, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				logger.Debug("error after parse", err.Error())
				errChan <- errors.Wrap(err, "Error read from body")
				return

			}

			logger.Debug("status", "success", "req", string(respData))

			r := csv.NewReader(strings.NewReader(string(respData)))

			for {
				record, err := r.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					errChan <- errors.Wrap(err, "read error csv")
					return
				}

				price, _ := strconv.ParseFloat(record[1], 64)
				product := repository.ProductPrice{Name: record[0], Price: float32(price)}

				err = s.repository.UpdatePrice(ctx, product)
				if err != nil {
					errChan <- errors.Wrap(err, fmt.Sprintf("err update %+v", product))
					return
				}
				return
			}
		}()
		for {
			select {
			case <-ctx.Done():
				tr.CancelRequest(req)
				return &FetchResponse{Status: "OK"}, nil
			case err := <-errChan:
				tr.CancelRequest(req)
				logger.Error("request", "canceled")
				return nil, errors.Wrap(err, "Request canceled")
			}

		}
	}

	return &FetchResponse{Status: "OK"}, nil
}
