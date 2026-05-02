package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/niuguy/langmate/llm"
)

const (
	defaultModel  = "openai-balanced"
	defaultLang   = "en"
	defaultHotkey = "cmd_ctrl_r"
)

type DaemonConfig struct {
	Model   string `json:"model"`
	Lang    string `json:"lang"`
	Preview bool   `json:"preview"`
	Hotkey  string `json:"hotkey"`
}

func DefaultDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Model:  defaultModel,
		Lang:   defaultLang,
		Hotkey: defaultHotkey,
	}
}

func LoadDaemonConfig() DaemonConfig {
	config := DefaultDaemonConfig()

	path, err := daemonConfigPath()
	if err != nil {
		return config
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultDaemonConfig()
	}

	config.normalize()
	return config
}

func SaveDaemonConfig(config DaemonConfig) error {
	config.normalize()

	path, err := daemonConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return os.WriteFile(path, data, 0o600)
}

func daemonConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "langmate", "config.json"), nil
}

func (c *DaemonConfig) normalize() {
	if c.Model == "" {
		c.Model = defaultModel
	}
	c.Model = llm.NormalizeModelPresetID(c.Model)
	if c.Lang == "" {
		c.Lang = defaultLang
	}
	if !isSupportedHotkey(c.Hotkey) {
		c.Hotkey = defaultHotkey
	}
}
