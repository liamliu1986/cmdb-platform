package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"cmdb-agent/core"
	"cmdb-agent/plugins"
	_ "cmdb-agent/plugins/aliyun"
	_ "cmdb-agent/plugins/tencent"
	_ "cmdb-agent/plugins/huawei"
	_ "cmdb-agent/plugins/aws"
)

func main() {
	cfg := core.LoadConfig()
	agent := core.NewAgent(cfg)

	// If no token, register first
	if cfg.Token == "" {
		fmt.Println("Registering agent...")
		if err := agent.Register(); err != nil {
			fmt.Println("Registration failed:", err)
			os.Exit(1)
		}
		fmt.Println("Registered successfully, token:", cfg.Token)
	}

	fmt.Println("Available plugins:")
	for name := range plugins.Registry {
		fmt.Println(" -", name)
	}

	fmt.Println("Agent started, sending heartbeats to", cfg.MasterURL)
	go agent.Run()

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
	agent.Stop()
}
