package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"time"
)

type ESClient struct {
	client *elastic.Client
}

type Post struct {
	Title, Body, User, Id string
	Timestamp             int64
	Likes                 int
}

var (
	ERR_ACKNOWLEDGEMENT_FAILED = errors.New("acknowledgement failed")
)

const (
	INDEX_NAME = "posts"
)

func NewESClient() (*ESClient, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}

	esc := &ESClient{client}
	if err = esc.createIndex(); err != nil {
		return nil, err
	}
	return esc, nil
}

func (esc *ESClient) createIndex() error {
	mappings := `{
		"mappings": {
			"properties": {
				"timestamp": {"type": "date"}
			}
		}	
	}`
	ctx := context.TODO()
	createIndex, err := esc.client.CreateIndex(INDEX_NAME).
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

func (esc *ESClient) Index(title, body, user, id string) error {
	ctx := context.TODO()
	_, err := esc.client.Index().
		Index(INDEX_NAME).
		BodyJson(Post{title, body, user, id, time.Now().Unix(), 0}).
		Do(ctx)
	return err
}

func (esc *ESClient) Delete(id string) error {
	ctx := context.TODO()
	query := elastic.NewTermQuery("id", id)
	_, err := esc.client.DeleteByQuery().
		Index(INDEX_NAME).
		Query(query).
		Do(ctx)
	return err
}

func (esc *ESClient) Search(term string, from, total int) ([]string, error) {
	ctx := context.TODO()
	query := elastic.NewMultiMatchQuery(term, "title", "body")
	searchResult, err := esc.client.Search().
		Index(INDEX_NAME).
		Query(query).
		Sort("timestamp", false).
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
	query := elastic.NewTermQuery("id", id)
	script := elastic.NewScript("ctx._source.likes = num").
		Param("num", likes)
	_, err := esc.client.UpdateByQuery().
		Index(INDEX_NAME).
		Script(script).
		Query(query).
		Do(ctx)
	return err
}
