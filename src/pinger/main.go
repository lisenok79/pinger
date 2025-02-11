package main

import (
	"context"
	"log"
	"os"
	containerinfo "pinger/containerInfo"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

func parseEnv() containerinfo.Env {
	var env containerinfo.Env
	networksEnv := os.Getenv("DOCKER_NETWORKS")
	if networksEnv == "" {
		log.Fatal("DOCKER_NETWORKS environment variable is required")
		os.Exit(1)
	}
	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		log.Fatal("BACKEND_URL environment variable is required")
		os.Exit(1)
	}
	networkList := strings.Split(networksEnv, ",")
	log.Printf("Monitoring networks: %v\n", networkList)
	env = containerinfo.Env{
		Networks: networkList,
		BackURL:  backendUrl,
	}
	return env
}

func main() {
	env := parseEnv()
	for {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal("Docker client init error: ", err)
		}
		if _, err := cli.Ping(context.Background()); err != nil {
			log.Fatal("Docker API connection error: ", err)
		}
		log.Println("Successfully connected to Docker API")
		containerinfo.CheckContainers(cli, env)
		time.Sleep(5 * time.Second)
		cli.Close()
	}
}
