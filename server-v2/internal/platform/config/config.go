package config

// Config 定义服务运行所需的最小配置集合。
type Config struct {
	ServiceName string `json:"service_name"`
	HTTPAddr    string `json:"http_addr"`
	LogLevel    string `json:"log_level"`
}

// ConfigOverlay 用于表达配置覆盖层，nil 表示该字段未提供。
type ConfigOverlay struct {
	ServiceName *string `json:"service_name"`
	HTTPAddr    *string `json:"http_addr"`
	LogLevel    *string `json:"log_level"`
}

// ApplyTo 将覆盖层应用到目标配置。
func (o ConfigOverlay) ApplyTo(dst *Config) {
	if o.ServiceName != nil {
		dst.ServiceName = *o.ServiceName
	}
	if o.HTTPAddr != nil {
		dst.HTTPAddr = *o.HTTPAddr
	}
	if o.LogLevel != nil {
		dst.LogLevel = *o.LogLevel
	}
}

// LoadOptions 用于控制配置加载来源。
type LoadOptions struct {
	FilePath string
	CLI      ConfigOverlay
}

// Default 返回默认配置。
func Default() Config {
	return Config{
		ServiceName: "server-v2",
		HTTPAddr:    ":8080",
		LogLevel:    "info",
	}
}
