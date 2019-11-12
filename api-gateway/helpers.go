package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

func parseJSON(body io.ReadCloser) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return data, err
	}
	err := json.Unmarshal(b, &data)
	return data, err
}

func get(addr string) (map[string]interface{}, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	data, err := parseJSON(resp.Body)
	if err != nil {
		retrun data, err
	}

	return data, err
}

func post(addr string, data interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{}, err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(respData))
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	data, err := parseJSON(resp.Body) 
	if err != nil {
		return data, err
	}

	if resp.StatusCode != http.StatusOK {
		return data, data["Error"]
	}

	return data, nil
}

func postProxy(addr string, body io.ReadCloser) (map[string]interface{}, error) {
	resp, err := http.Post(addr, "application/json", body)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	data, err := parseJSON(resp.Body) 
	if err != nil {
		return data, err
	}

	if resp.StatusCode != http.StatusOK {
		return data, data["Error"]
	}

	return data, nil
}

func bindJSONAndAdd(ctx *gin.Context, other map[string]interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := ctx.ShouldBindJSON(&m)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	for k, v := range other {
		m[k] = v
	}

	return m, nil
}