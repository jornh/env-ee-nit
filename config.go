package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

// LoadConfig reads the TOML file into the Config struct
func LoadConfig(path string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// AddEnvIfMissing adds the environment to the list if it doesn't exist
func (c *Config) AddEnvIfMissing(env string) {
	for _, e := range c.Envs {
		if e == env {
			return
		}
	}
	c.Envs = append(c.Envs, env)
}

// SaveConfig writes the current struct back to the TOML file
func (c *Config) SaveConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(c)
}

