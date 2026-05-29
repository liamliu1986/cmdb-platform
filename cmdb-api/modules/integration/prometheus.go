package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type PrometheusClient struct {
	BaseURL string
}

func NewPrometheusClient(baseURL string) *PrometheusClient {
	if baseURL == "" {
		baseURL = "http://localhost:9090"
	}
	return &PrometheusClient{BaseURL: baseURL}
}

// QueryInstant performs an instant query against Prometheus
func (p *PrometheusClient) QueryInstant(query string) (interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1/query?query=%s", p.BaseURL, url.QueryEscape(query))
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("prometheus query failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)
	return result, nil
}

// QueryRange performs a range query against Prometheus
func (p *PrometheusClient) QueryRange(query string, start, end, step string) (interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		p.BaseURL, url.QueryEscape(query), start, end, step)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("prometheus range query failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)
	return result, nil
}
