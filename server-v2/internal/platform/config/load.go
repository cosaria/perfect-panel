package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Load 按 file -> env -> CLI 顺序合并配置，最终优先级为 CLI > env > file。
func Load(opts LoadOptions) (Config, error) {
	cfg := Default()

	if opts.FilePath != "" {
		fileOverlay, err := loadFromFile(opts.FilePath)
		if err != nil {
			return Config{}, err
		}
		fileOverlay.ApplyTo(&cfg)
	}

	loadFromEnv().ApplyTo(&cfg)
	opts.CLI.ApplyTo(&cfg)

	return cfg, nil
}

func loadFromFile(path string) (ConfigOverlay, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ConfigOverlay{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var overlay ConfigOverlay
	if err := json.Unmarshal(content, &overlay); err != nil {
		return ConfigOverlay{}, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return overlay, nil
}

func loadFromEnv() ConfigOverlay {
	return ConfigOverlay{
		ServiceName: lookupEnvPtr("PPANEL_SERVICE_NAME"),
		HTTPAddr:    lookupEnvPtr("PPANEL_HTTP_ADDR"),
		LogLevel:    lookupEnvPtr("PPANEL_LOG_LEVEL"),
	}
}

func lookupEnvPtr(key string) *string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}
	return &value
}
