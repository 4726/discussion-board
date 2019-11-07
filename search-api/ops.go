package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type ESClient struct {
	client    *elastic.Client
	indexName string
}

type Post struct {
	Title, Body, User, Id string
	LastUpdate            int64
	Likes                 int
}

var (
	ERR_ACKNOWLEDGEMENT_FAILED = errors.New("acknowledgement failed")
)

func NewESClient(indexName string) (*ESClient, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	esc := &ESClient{client, indexName}
	if err = esc.createIndex(); err != nil {
		return nil, err
	}
	return esc, nil
}

func (esc *ESClient) createIndex() error {
	mappings := `{
		"mappings": {
			"properties": {
				"Title": {"type": "text"},
				"Body": {"type": "text"},
				"User": {"type": "keyword"},
				"Id": {"type": "keyword"},
				"LastUpdate": {"type": "integer"},
				"Likes": {"type": "integer"}
			}
		}	
	}`
	ctx := context.TODO()
	createIndex, err := esc.client.CreateIndex(esc.indexName).
		BodyString(mappings).
		Do(ctx)
	if err != nil {
		if e, ok := err.(*elastic.Error); ok {
			if e.Details.Type == "resource_already_exists_exception" {
				return nil
			}
		}
		return err
	}
	if !createIndex.Acknowledged {
		return ERR_ACKNOWLEDGEMENT_FAILED
	}
	return nil
}

//is not immediately available for search,
//can set refresh to wait until index completes before returning
func (esc *ESClient) Index(post Post) error {
	ctx := context.TODO()
	_, err := esc.client.Index().
		Index(esc.indexName).
		BodyJson(post).
		Do(ctx)
	return err
}

func (esc *ESClient) Search(term string, from, total int) ([]string, error) {
	ctx := context.TODO()
	query := elastic.NewMultiMatchQuery(term, "Title", "Body")
	searchResult, err := esc.client.Search().
		Index(esc.indexName).
		Query(query).
		SortBy(elastic.NewFieldSort("LastUpdate").Asc()).
		From(from).Size(total).
		Do(ctx)
	if err != nil {
		return []string{}, err
	}
	results := []string{}
	for _, hit := range searchResult.Hits.Hits {
		var p Post
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			return []string{}, err
		}
		results = append(results, p.Id)
	}

	return results, nil
}

func (esc *ESClient) UpdateLikes(id string, likes int) error {
	ctx := context.TODO()
	query := elastic.NewTermQuery("Id", id)
	script := elastic.NewScript(fmt.Sprintf("ctx._source.Likes = %v", likes))
	_, err := esc.client.UpdateByQuery().
		Index(esc.indexName).
		Query(query).
		Script(script).
		Do(ctx)
	return err
}

func (esc *ESClient) Delete(id string) error {
	ctx := context.TODO()
	query := elastic.NewTermQuery("Id", id)
	_, err := esc.client.DeleteByQuery().
		Index(esc.indexName).
		Query(query).
		Do(ctx)
	return err
}

func (esc *ESClient) UpdateLastUpdate(id string, lastUpdate int64) error {
	ctx := context.TODO()
	query := elastic.NewTermQuery("Id", id)
	script := elastic.NewScript(fmt.Sprintf("ctx._source.LastUpdate = %v", lastUpdate))
	_, err := esc.client.UpdateByQuery().
		Index(esc.indexName).
		Query(query).
		Script(script).
		Do(ctx)
	return err
}