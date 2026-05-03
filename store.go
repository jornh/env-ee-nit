package main

type ConfigStore interface {
    Load() (*Config, string, error)
    Save(*Config, string) error
}
