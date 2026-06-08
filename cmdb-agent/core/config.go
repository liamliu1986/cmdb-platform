package core

import (
	"os"
)

type Config struct {
	MasterURL string
	Token     string
	AgentID   uint
	Name      string
	IP        string
	OS        string
	Arch      string
	Version   string
}

func LoadConfig() *Config {
	cfg := &Config{
		MasterURL: getEnv("MASTER_URL", "http://localhost:8080"),
		Token:     getEnv("AGENT_TOKEN", ""),
		Name:      getEnv("AGENT_NAME", ""),
		IP:        getEnv("AGENT_IP", ""),
		OS:        getEnv("AGENT_OS", ""),
		Arch:      getEnv("AGENT_ARCH", ""),
		Version:   "1.0.0",
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
