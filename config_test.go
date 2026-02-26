package main

import (
	"os"
	"testing"
)

func TestAddEnvIfMissing(t *testing.T) {
	cfg := &Config{
		Envs: []string{"dev", "prod"},
	}

	// Test adding a new one
	cfg.AddEnvIfMissing("staging")
	if len(cfg.Envs) != 3 || cfg.Envs[2] != "staging" {
		t.Errorf("Expected staging to be added, got %v", cfg.Envs)
	}

	// Test adding a duplicate (should do nothing)
	cfg.AddEnvIfMissing("dev")
	if len(cfg.Envs) != 3 {
		t.Errorf("Expected length to stay 3, got %d", len(cfg.Envs))
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpFile := "test_versions.toml"
	defer os.Remove(tmpFile)

	initial := &Config{
		Envs:      []string{"dev"},
		GithubOrg: "my-org",
	}

	// Save
	err := initial.SaveConfig(tmpFile)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.GithubOrg != "my-org" || loaded.Envs[0] != "dev" {
		t.Errorf("Data mismatch after reload. Got: %+v", loaded)
	}
}
