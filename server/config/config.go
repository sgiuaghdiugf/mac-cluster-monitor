package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 总配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Monitor  MonitorConfig  `yaml:"monitor"`
	Profit   ProfitConfig   `yaml:"profit"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	OfflineThreshold int `yaml:"offline_threshold"`
	CheckInterval    int `yaml:"check_interval"`
}

// ProfitConfig 盈利模拟配置
type ProfitConfig struct {
	InitialValue float64 `yaml:"initial_value"` // 初始盈利值
	Increment    float64 `yaml:"increment"`     // 每次增加的盈利
	Interval     int     `yaml:"interval"`      // 增加间隔(秒)
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load 从YAML文件加载配置
func Load(configPath string) (*Config, error) {
	// 默认配置
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:         "localhost",
			Port:         "3306",
			User:         "root",
			Password:     "",
			Name:         "mac_cluster_monitor",
			MaxOpenConns: 25,
			MaxIdleConns: 10,
		},
		JWT: JWTConfig{
			Secret:      "default-secret-key",
			ExpireHours: 24,
		},
		Monitor: MonitorConfig{
			OfflineThreshold: 10,
			CheckInterval:    5,
		},
		Profit: ProfitConfig{
			InitialValue: 0,
			Increment:    0.5,
			Interval:     5,
		},
		Log: LogConfig{
			Level:  "debug",
			Format: "text",
		},
	}

	// 如果指定了配置文件路径,从文件加载
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	return cfg, nil
}

// LoadFromEnv 从环境变量获取配置文件路径并加载
func LoadFromEnv() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认查找当前目录下的 config.yaml
		if _, err := os.Stat("config.yaml"); err == nil {
			configPath = "config.yaml"
		}
	}

	cfg, err := Load(configPath)
	if err != nil {
		return nil, err
	}

	// 环境变量覆盖（Docker部署时使用）
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		cfg.Database.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.Name = v
	}
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}

	return cfg, nil
}

// DSN 生成MySQL连接字符串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}
