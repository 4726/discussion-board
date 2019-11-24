package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
)

type Resp struct {
	Data       map[string]interface{}
	StatusCode int
}

type RespArray struct {
	Data       []interface{}
	StatusCode int
}


func parseJSON(body io.ReadCloser) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}

func parseJSONArray(body io.ReadCloser) ([]interface{}, error) {
	data := []interface{}{}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	return data, err
}

func get(addr string) (Resp, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return Resp{map[string]interface{}{}, 0}, err
	}
	defer resp.Body.Close()

	data, err := parseJSON(resp.Body)
	if err != nil {
		return Resp{data, resp.StatusCode}, err
	}

	return Resp{data, resp.StatusCode}, err
}

func post(addr string, data interface{}) (Resp, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return Resp{map[string]interface{}{}, 0}, err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return Resp{map[string]interface{}{}, 0}, err
	}
	defer resp.Body.Close()

	respData, err := parseJSON(resp.Body)
	if err != nil {
		return Resp{respData, resp.StatusCode}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Resp{respData, resp.StatusCode}, fmt.Errorf(respData["Error"].(string))
	}

	return Resp{respData, resp.StatusCode}, nil
}

func postArray(addr string, data interface{}) (RespArray, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return RespArray{[]interface{}{}, 0}, err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return RespArray{[]interface{}{}, 0}, err
	}
	defer resp.Body.Close()

	respData, err := parseJSONArray(resp.Body)
	if err != nil {
		return RespArray{respData, resp.StatusCode}, err
	}

	if resp.StatusCode != http.StatusOK {
		return RespArray{respData, resp.StatusCode}, nil
	}

	return RespArray{respData, resp.StatusCode}, nil
}

func postProxy(addr string, body io.ReadCloser) (Resp, error) {
	defer func() {
		if body != nil {
			body.Close()
		}
	}()
	resp, err := http.Post(addr, "application/json", body)
	if err != nil {
		return Resp{map[string]interface{}{}, 0}, err
	}
	defer resp.Body.Close()

	data, err := parseJSON(resp.Body)
	if err != nil {
		return Resp{map[string]interface{}{}, resp.StatusCode}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Resp{data, resp.StatusCode}, fmt.Errorf(data["Error"].(string))
	}

	return Resp{data, resp.StatusCode}, nil
}

func bindJSONAndAdd(ctx *gin.Context, other map[string]interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := ctx.ShouldBindJSON(&m)
	if err != nil {
		return m, nil
	}

	for k, v := range other {
		m[k] = v
	}

	return m, nil
}
