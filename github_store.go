package main

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "github.com/BurntSushi/toml"
)

type GitHubStore struct {
    Host   string
    Token  string
    Owner  string
    Repo   string
    Path   string
    Client *http.Client
}

func NewGitHubStore(spec string) (*GitHubStore, error) {
    // spec format: org/repo/path/to/file.toml
    parts := strings.SplitN(spec, "/", 3)
    if len(parts) != 3 {
        return nil, fmt.Errorf("invalid --github spec, expected org/repo/path")
    }

    host := os.Getenv("GH_HOST")
    if host == "" {
        host = "api.github.com"
    }

    token := os.Getenv("GH_TOKEN")
    if token == "" {
        return nil, fmt.Errorf("GH_TOKEN is required")
    }

    return &GitHubStore{
        Host:   host,
        Token:  token,
        Owner:  parts[0],
        Repo:   parts[1],
        Path:   parts[2],
        Client: http.DefaultClient,
    }, nil
}

func (g *GitHubStore) apiURL() string {
    return fmt.Sprintf("https://%s/repos/%s/%s/contents/%s",
        g.Host, g.Owner, g.Repo, g.Path)
}

func (g *GitHubStore) Load() (*Config, string, error) {
    req, _ := http.NewRequest("GET", g.apiURL(), nil)
    req.Header.Set("Authorization", "token "+g.Token)

    resp, err := g.Client.Do(req)
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
        // No file yet → empty config
        return &Config{}, "", nil
    }

    if resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        return nil, "", fmt.Errorf("GitHub GET failed: %s", string(body))
    }

    var data struct {
        Content string `json:"content"`
        SHA     string `json:"sha"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, "", err
    }

    raw, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(data.Content, "\n", ""))
    if err != nil {
        return nil, "", err
    }

    var cfg Config
    if _, err := toml.Decode(string(raw), &cfg); err != nil {
        return nil, "", err
    }

    return &cfg, data.SHA, nil
}

func (g *GitHubStore) Save(cfg *Config, sha string) error {
    var buf bytes.Buffer
    if err := toml.NewEncoder(&buf).Encode(cfg); err != nil {
        return err
    }

    body := map[string]string{
        "message": "Update versions.toml via env-ee-nit",
        "content": base64.StdEncoding.EncodeToString(buf.Bytes()),
    }
    if sha != "" {
        body["sha"] = sha
    }

    b, _ := json.Marshal(body)

    req, _ := http.NewRequest("PUT", g.apiURL(), bytes.NewReader(b))
    req.Header.Set("Authorization", "token "+g.Token)
    req.Header.Set("Content-Type", "application/json")

    resp, err := g.Client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        out, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("GitHub PUT failed: %s", string(out))
    }

    return nil
}
