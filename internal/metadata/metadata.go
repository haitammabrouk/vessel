package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const MetadataPath = "/var/lib/vessel/containers"
var configPath string

type Config struct {
	ID         string  `json:"ID"`
	Created    string  `json:"Created"`
	ConfigPath string  `json:"ConfigPath"`
}

//NewConfig set up a config dir for the given container id
func InitContainerConfig(id, created string) error {
	cfg := &Config{
		ID: id,
		Created: created,
	}

	if err := setUpMetadataPath(); err != nil {
		return err
	}

	if err := cfg.setUpConfigDir(); err != nil {
		return err
	}

	if err := cfg.writeConfigInfos(); err != nil {
		return err
	}
	return nil
}

func setUpMetadataPath() error {
	if err := os.MkdirAll(MetadataPath, 0700); err != nil {
		return fmt.Errorf("create metadata dir for containers: %w", err)
	}
	return nil
}

func (cfg *Config) setUpConfigDir() error {
	configPath = filepath.Join(MetadataPath, cfg.ID)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("create config dir for container: %w", err)
	}
	return nil
}

func (cfg *Config) writeConfigInfos() error {
	cfg.ConfigPath = filepath.Join(configPath, "config.json")
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode container config to json: %w", err)
	}
	
	if err := os.WriteFile(cfg.ConfigPath, data, 0644); err != nil {
		return fmt.Errorf("write to config file of container :%w", err)
	}
	return nil
}