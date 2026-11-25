package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.WhatsAppEnabled {
		t.Errorf("WhatsAppEnabled debe estar habilitado por defecto")
	}
	if cfg.HealthPort == "" || cfg.Model == "" || cfg.SystemPrompt == "" {
		t.Errorf("DefaultConfig contiene campos criticos vacios: %+v", cfg)
	}
}

func TestLoadConfigCreatesDefaultAndNormalizes(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig fallo creando config por defecto: %v", err)
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("config.json no fue creado en disco: %v", err)
	}
	if !cfg.WhatsAppEnabled {
		t.Errorf("WhatsAppEnabled deberia ser true")
	}

	if err := os.WriteFile(cfgPath, []byte(`{"whatsapp_enabled": false, "max_iterations": 0}`), 0600); err != nil {
		t.Fatalf("error escribiendo config parcial: %v", err)
	}
	normalized, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig fallo leyendo config parcial: %v", err)
	}
	if !normalized.WhatsAppEnabled || normalized.MaxIterations != 15 || normalized.HealthPort != "8055" {
		t.Errorf("LoadConfig no normalizo correctamente: %+v", normalized)
	}
}

func TestConfigSaveAndSetOwner(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	cfg := DefaultConfig()

	if err := cfg.SetOwner("5491112345678@s.whatsapp.net", cfgPath); err != nil {
		t.Fatalf("SetOwner fallo: %v", err)
	}

	reloaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadConfig tras Save fallo: %v", err)
	}
	if reloaded.OwnerJID != "5491112345678@s.whatsapp.net" {
		t.Errorf("OwnerJID inesperado: %q", reloaded.OwnerJID)
	}
	if len(reloaded.AllowedUsers) != 1 || reloaded.AllowedUsers[0] != "5491112345678@s.whatsapp.net" {
		t.Errorf("AllowedUsers inesperado: %+v", reloaded.AllowedUsers)
	}
}
