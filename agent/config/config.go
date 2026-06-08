package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config Agent配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	Device DeviceConfig `yaml:"device"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	URL               string `yaml:"url"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
	Timeout           int    `yaml:"timeout"`
}

// DeviceConfig 设备配置
type DeviceConfig struct {
	Name      string `yaml:"name"`
	ID        string `yaml:"id"`
	Interface string `yaml:"interface"`
	Region    string `yaml:"region"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// Load 从YAML文件加载配置
func Load(configPath string) (*Config, error) {
	// 默认配置
	cfg := &Config{
		Server: ServerConfig{
			URL:               "http://localhost:8080",
			HeartbeatInterval: 5,
			Timeout:           10,
		},
		Device: DeviceConfig{
			Name:      "",
			ID:        "",
			Interface: "",
			Region:    "北京", // 默认地区，根据实际情况修改
		},
		Log: LogConfig{
			Level: "info",
			File:  "/var/log/mac-monitor-agent.log",
		},
	}

	// 如果指定了配置文件路径，从文件加载
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 环境变量覆盖配置文件
	if url := os.Getenv("MONITOR_SERVER"); url != "" {
		cfg.Server.URL = url
	}
	if interval := os.Getenv("MONITOR_INTERVAL"); interval != "" {
		var intVal int
		if _, err := fmt.Sscanf(interval, "%d", &intVal); err == nil {
			cfg.Server.HeartbeatInterval = intVal
		}
	}

	return cfg, nil
}

// LoadFromEnv 从环境变量获取配置文件路径并加载
func LoadFromEnv() (*Config, error) {
	configPath := os.Getenv("AGENT_CONFIG_PATH")
	if configPath == "" {
		// 默认查找当前目录下的 config.yaml
		if _, err := os.Stat("config.yaml"); err == nil {
			configPath = "config.yaml"
		}
	}

	return Load(configPath)
}
