package esutils

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"
	"github.com/faisalyudiansah/auth-service-template/pkg/logger"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ESUtils interface {
	Create(ctx context.Context, index string, id string, document any) error
	Get(ctx context.Context, index string, id string) (map[string]any, error)
	Update(ctx context.Context, index string, id string, document any) error
	Delete(ctx context.Context, index string, id string) error
	Search(ctx context.Context, index string, query any) (*esapi.Response, error)
}

type esUtils struct {
	client *elasticsearch.Client
}

func NewESUtils(cfg *config.ESConfig) *esUtils {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
	})
	if err != nil {
		log.Fatalf("error creating es client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("error getting es response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("error es response: %s", res.String())
	}

	var r map[string]any
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("error parsing the es response body: %s", err)
	}

	logger.Log.Infof("es client: %s", elasticsearch.Version)
	logger.Log.Infof("es server: %s", r["version"].(map[string]any)["number"])

	return &esUtils{
		client: es,
	}
}

func (e *esUtils) Create(ctx context.Context, index string, id string, document any) error {
	body, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	return nil
}

func (e *esUtils) Get(ctx context.Context, index string, id string) (map[string]any, error) {
	req := esapi.GetRequest{
		Index:      index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, err
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (e *esUtils) Update(ctx context.Context, index string, id string, document any) error {
	body, err := json.Marshal(document)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	return nil
}

func (e *esUtils) Delete(ctx context.Context, index string, id string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	return nil
}

func (e *esUtils) Search(ctx context.Context, index string, query interface{}) (*esapi.Response, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{index},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, err
	}

	return res, nil
}
