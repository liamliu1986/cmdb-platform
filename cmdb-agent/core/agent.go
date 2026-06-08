package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type Agent struct {
	cfg      *Config
	client   *http.Client
	stopChan chan struct{}
}

func NewAgent(cfg *Config) *Agent {
	return &Agent{
		cfg:      cfg,
		client:   &http.Client{Timeout: 30 * time.Second},
		stopChan: make(chan struct{}),
	}
}

func (a *Agent) getLocalIP() string {
	// Simplified: just return empty, user can set via env
	return a.cfg.IP
}

func (a *Agent) Register() error {
	reqBody := map[string]interface{}{
		"name":    a.cfg.Name,
		"ip":      a.getLocalIP(),
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": a.cfg.Version,
	}
	body, _ := json.Marshal(reqBody)
	resp, err := a.client.Post(
		a.cfg.MasterURL+"/api/v1/discovery/agents/register",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("register failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Code int `json:"code"`
		Data struct {
			AgentID uint   `json:"agent_id"`
			Token   string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Code != 0 {
		return fmt.Errorf("register failed with code: %d", result.Code)
	}
	a.cfg.AgentID = result.Data.AgentID
	a.cfg.Token = result.Data.Token
	return nil
}

func (a *Agent) Heartbeat() error {
	reqBody := map[string]string{"status": "running"}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(
		"POST",
		a.cfg.MasterURL+"/api/v1/discovery/agents/heartbeat",
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Agent-Token", a.cfg.Token)

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *Agent) Run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// First heartbeat immediately
	if err := a.Heartbeat(); err != nil {
		fmt.Println("Initial heartbeat failed:", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := a.Heartbeat(); err != nil {
				fmt.Println("Heartbeat failed:", err)
			}
		case <-a.stopChan:
			return
		}
	}
}

func (a *Agent) Stop() {
	close(a.stopChan)
}
