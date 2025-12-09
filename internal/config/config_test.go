package config

import (
	"os"
	"testing"
)

func TestLoadConfig_YAML(t *testing.T) {
	content := `
users:
  - name: "test-user"
    token: "test-token"
listen_addr: ":9102"
metrics_path: "/test-metrics"
poll_interval: 30
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(cfg.Users))
	}

	if cfg.Users[0].Name != "test-user" {
		t.Errorf("Expected user name 'test-user', got '%s'", cfg.Users[0].Name)
	}

	if cfg.Users[0].Token != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", cfg.Users[0].Token)
	}

	if cfg.ListenAddr != ":9102" {
		t.Errorf("Expected listen_addr ':9102', got '%s'", cfg.ListenAddr)
	}

	if cfg.MetricsPath != "/test-metrics" {
		t.Errorf("Expected metrics_path '/test-metrics', got '%s'", cfg.MetricsPath)
	}

	if cfg.PollInterval != 30 {
		t.Errorf("Expected poll_interval 30, got %d", cfg.PollInterval)
	}
}

func TestLoadConfig_TOML(t *testing.T) {
	content := `
listen_addr = ":9103"
metrics_path = "/metrics"
poll_interval = 45

[[users]]
name = "user1"
token = "token1"

[[users]]
name = "user2"
token = "token2"
`
	tmpfile, err := os.CreateTemp("", "config-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(cfg.Users))
	}

	if cfg.Users[0].Name != "user1" {
		t.Errorf("Expected first user name 'user1', got '%s'", cfg.Users[0].Name)
	}

	if cfg.Users[1].Name != "user2" {
		t.Errorf("Expected second user name 'user2', got '%s'", cfg.Users[1].Name)
	}
}

func TestLoadConfig_HCL(t *testing.T) {
	content := `
listen_addr = ":9104"
metrics_path = "/metrics"
poll_interval = 90

user {
  name = "hcl-user"
  token = "hcl-token"
}
`
	tmpfile, err := os.CreateTemp("", "config-*.hcl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(cfg.Users))
	}

	if cfg.Users[0].Name != "hcl-user" {
		t.Errorf("Expected user name 'hcl-user', got '%s'", cfg.Users[0].Name)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	content := `
users:
  - name: "test-user"
    token: "test-token"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ListenAddr != ":9101" {
		t.Errorf("Expected default listen_addr ':9101', got '%s'", cfg.ListenAddr)
	}

	if cfg.MetricsPath != "/metrics" {
		t.Errorf("Expected default metrics_path '/metrics', got '%s'", cfg.MetricsPath)
	}

	if cfg.PollInterval != 60 {
		t.Errorf("Expected default poll_interval 60, got %d", cfg.PollInterval)
	}
}

func TestLoadConfig_NoUsers(t *testing.T) {
	content := `
listen_addr: ":9101"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for config with no users, got nil")
	}
}

func TestLoadConfig_MissingUserName(t *testing.T) {
	content := `
users:
  - token: "test-token"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for user with missing name, got nil")
	}
}

func TestLoadConfig_MissingToken(t *testing.T) {
	content := `
users:
  - name: "test-user"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for user with missing token, got nil")
	}
}

func TestLoadConfig_UnsupportedFormat(t *testing.T) {
	content := `invalid content`
	tmpfile, err := os.CreateTemp("", "config-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for unsupported file format, got nil")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	content := `
users:
  - name: "test-user"
    token: "test-token"
  invalid yaml syntax here
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	_, err = LoadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadConfig_MultipleUsers(t *testing.T) {
	content := `
users:
  - name: "user1"
    token: "token1"
  - name: "user2"
    token: "token2"
  - name: "user3"
    token: "token3"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(cfg.Users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(cfg.Users))
	}

	expectedUsers := []struct {
		name  string
		token string
	}{
		{"user1", "token1"},
		{"user2", "token2"},
		{"user3", "token3"},
	}

	for i, expected := range expectedUsers {
		if cfg.Users[i].Name != expected.name {
			t.Errorf("User %d: expected name '%s', got '%s'", i, expected.name, cfg.Users[i].Name)
		}
		if cfg.Users[i].Token != expected.token {
			t.Errorf("User %d: expected token '%s', got '%s'", i, expected.token, cfg.Users[i].Token)
		}
	}
}

func TestLoadConfig_YMLExtension(t *testing.T) {
	content := `
users:
  - name: "test-user"
    token: "test-token"
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config with .yml extension: %v", err)
	}

	if len(cfg.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(cfg.Users))
	}
}
