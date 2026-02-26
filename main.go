package main

import (
	"flag"
	"fmt"
	"log"
	"k8s.io/client-go/kubernetes" // This provides the Interface type
)

type Config struct {
        Envs            []string  `toml:"envs"`
        GithubOrg       string    `toml:"github_org"`
        GitTagTransform string    `toml:"git_tag_transform"`
        Versions        []Version `toml:"versions"`
}

type Version struct {
        App     string `toml:"app"`
        Env     string `toml:"env"`
        Version string `toml:"version"`
}

func main() {
	envPtr := flag.String("env", "", "The environment/namespace to refresh")
	localDir := flag.String("local", "", "Path to local k8s yaml files (optional)")
	configPath := flag.String("versions", "versions.toml", "Path to output envee versions file")
	flag.Parse()

	if *envPtr == "" {
		log.Fatal("Error: --env parameter is required")
	}

	// Initialize the appropriate client
	var client kubernetes.Interface
	var err error

	if *localDir != "" {
		fmt.Printf("Using local YAML files from: %s\n", *localDir)
		client, err = GetK8sFilesClient(*localDir)
	} else {
		client, err = GetK8sClient()
	}

	if err != nil {
		log.Fatalf("Client Error: %v", err)
	}

	// 1. Load existing config
	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Update environment list
	RunRefresh(cfg, *envPtr) // Clears old data for this env and adds to Envs list
	// TODO: still needed?
	// cfg.AddEnvIfMissing(*envPtr)

	fmt.Printf("Refreshing environment: %s\n", *envPtr)

	// Step 3 & 4: Fetch from K8s and merge
	k8sVersions, err := FetchVersionsFromNamespace(client, *envPtr)
	if err != nil {
		log.Fatalf("Error fetching versions: %v", err)
	}
	cfg.Versions = append(cfg.Versions, k8sVersions...)


	// 5. Save back to file
	if err := cfg.SaveConfig(*configPath); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	// Accessing data
        fmt.Printf("Github Org: %s\n", cfg.GithubOrg)
        for _, v := range cfg.Versions {
                fmt.Printf("App: %s | Env: %s | Version: %s\n", v.App, v.Env, v.Version)
        }
}
