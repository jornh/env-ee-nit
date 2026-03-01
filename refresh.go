package main

import (
	"sort"
)

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

// SortVersions organizes the versions slice by App name, then by the index of Env in Envs list
func (cfg *Config) SortVersions() {
	// Create a map for O(1) lookup of environment priority
	envPriority := make(map[string]int)
	for i, env := range cfg.Envs {
		envPriority[env] = i
	}

	sort.Slice(cfg.Versions, func(i, j int) bool {
		v1 := cfg.Versions[i]
		v2 := cfg.Versions[j]

		// Primary sort: App Name (alphabetical)
		if v1.App != v2.App {
			return v1.App < v2.App
		}

		// Secondary sort: Env (based on order in cfg.Envs)
		p1, ok1 := envPriority[v1.Env]
		p2, ok2 := envPriority[v2.Env]

		// If an env isn't in the list for some reason, push it to the end
		if !ok1 { p1 = len(cfg.Envs) }
		if !ok2 { p2 = len(cfg.Envs) }

		return p1 < p2
	})
}
