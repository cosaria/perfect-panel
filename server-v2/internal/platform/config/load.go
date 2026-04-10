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
		fileCfg, err := loadFromFile(opts.FilePath)
		if err != nil {
			return Config{}, err
		}
		mergeConfig(&cfg, fileCfg)
	}

	mergeConfig(&cfg, loadFromEnv())
	mergeConfig(&cfg, opts.CLI)

	return cfg, nil
}

func loadFromFile(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return Config{}, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return cfg, nil
}

func loadFromEnv() Config {
	return Config{
		ServiceName: getEnv("PPANEL_SERVICE_NAME"),
		HTTPAddr:    getEnv("PPANEL_HTTP_ADDR"),
		LogLevel:    getEnv("PPANEL_LOG_LEVEL"),
	}
}

func getEnv(key string) string {
	v, _ := os.LookupEnv(key)
	return v
}

func mergeConfig(dst *Config, src Config) {
	if src.ServiceName != "" {
		dst.ServiceName = src.ServiceName
	}
	if src.HTTPAddr != "" {
		dst.HTTPAddr = src.HTTPAddr
	}
	if src.LogLevel != "" {
		dst.LogLevel = src.LogLevel
	}
}
