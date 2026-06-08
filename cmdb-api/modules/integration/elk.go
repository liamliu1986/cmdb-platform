package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ELKClient struct {
	BaseURL string
	Index   string
}

func NewELKClient(baseURL, index string) *ELKClient {
	if baseURL == "" {
		baseURL = "http://localhost:9200"
	}
	if index == "" {
		index = "logstash-*"
	}
	return &ELKClient{BaseURL: baseURL, Index: index}
}

// SearchLogs searches logs by hostname/IP
func (e *ELKClient) SearchLogs(query string, size int) (interface{}, error) {
	if size <= 0 {
		size = 50
	}
	searchBody := map[string]interface{}{
		"size": size,
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": query,
			},
		},
		"sort": []map[string]interface{}{
			{"@timestamp": map[string]string{"order": "desc"}},
		},
	}
	body, _ := json.Marshal(searchBody)
	endpoint := fmt.Sprintf("%s/%s/_search", e.BaseURL, e.Index)
	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ELK search failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}
