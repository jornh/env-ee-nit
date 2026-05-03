package main

import (
    "os"

    "github.com/BurntSushi/toml"
)

type FileStore struct {
    Path string
}

func (f *FileStore) Load() (*Config, string, error) {
    var cfg Config
    if _, err := os.Stat(f.Path); os.IsNotExist(err) {
        return &cfg, "", nil
    }
    _, err := toml.DecodeFile(f.Path, &cfg)
    return &cfg, "", err
}

func (f *FileStore) Save(cfg *Config, _ string) error {
    file, err := os.Create(f.Path)
    if err != nil {
        return err
    }
    defer file.Close()
    return toml.NewEncoder(file).Encode(cfg)
}
