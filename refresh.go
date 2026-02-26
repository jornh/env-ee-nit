package main

//import (
//	"fmt"
//)

// RunRefresh performs the logic for Step 1 & 2
func RunRefresh(cfg *Config, targetEnv string) {
	// Step 1: Validate/Add the environment to the global list
	exists := false
	for _, e := range cfg.Envs {
		if e == targetEnv {
			exists = true
			break
		}
	}
	if !exists {
		cfg.Envs = append(cfg.Envs, targetEnv)
	}

	// Step 2: Clear old versions for this specific env (to be replaced by K8s data)
	// We filter the slice to keep only entries for OTHER environments
	newVersions := make([]Version, 0)
	for _, v := range cfg.Versions {
		if v.Env != targetEnv {
			newVersions = append(newVersions, v)
		}
	}
	cfg.Versions = newVersions
}

