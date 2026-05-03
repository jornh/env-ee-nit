package main

import (
	"os"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	tmpFile := "test_versions.toml"
	var store ConfigStore
	store = &FileStore{Path: tmpFile}
	defer os.Remove(tmpFile)

	initial := &Config{
		Envs:      []string{"dev"},
		GithubOrg: "my-org",
	}

	// Save
	err := store.Save(initial, "")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, _, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.GithubOrg != "my-org" || loaded.Envs[0] != "dev" {
		t.Errorf("Data mismatch after reload. Got: %+v", loaded)
	}
}
