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
	Title, Body string
	Id          uint64
	UserID      uint64
	LastUpdate  int64
	Likes       int64
}

var (
	ERR_ACKNOWLEDGEMENT_FAILED = errors.New("acknowledgement failed")
)

func NewESClient(indexName, addr string) (*ESClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(addr),
	)
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
				"UserID": {"type": "integer"},
				"Id": {"type": "integer"},
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

func (esc *ESClient) Search(term string, from, total uint64) ([]uint64, error) {
	ctx := context.TODO()
	query := elastic.NewMultiMatchQuery(term, "Title", "Body")
	searchResult, err := esc.client.Search().
		Index(esc.indexName).
		Query(query).
		SortBy(elastic.NewFieldSort("LastUpdate").Asc()).
		From(int(from)).Size(int(total)).
		Do(ctx)
	if err != nil {
		return []uint64{}, err
	}
	results := []uint64{}
	for _, hit := range searchResult.Hits.Hits {
		var p Post
		if err := json.Unmarshal(hit.Source, &p); err != nil {
			return []uint64{}, err
		}
		results = append(results, p.Id)
	}

	return results, nil
}

func (esc *ESClient) UpdateLikes(id uint64, likes int64) error {
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

func (esc *ESClient) Delete(id uint64) error {
	ctx := context.TODO()
	query := elastic.NewTermQuery("Id", id)
	_, err := esc.client.DeleteByQuery().
		Index(esc.indexName).
		Query(query).
		Do(ctx)
	return err
}

func (esc *ESClient) UpdateLastUpdate(id uint64, lastUpdate int64) error {
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
