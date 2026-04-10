package config

// Config 定义服务运行所需的最小配置集合。
type Config struct {
	ServiceName string `json:"service_name"`
	HTTPAddr    string `json:"http_addr"`
	LogLevel    string `json:"log_level"`
}

// LoadOptions 用于控制配置加载来源。
type LoadOptions struct {
	FilePath string
	CLI      Config
}

// Default 返回默认配置。
func Default() Config {
	return Config{
		ServiceName: "server-v2",
		HTTPAddr:    ":8080",
		LogLevel:    "info",
	}
}
