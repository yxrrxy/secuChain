package config

import (
	"blockSBOM/pkg/utils"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port    int    `yaml:"port"`
		Mode    string `yaml:"mode"`
		Address string `yaml:"address"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host" default:"localhost"`
		Port     int    `yaml:"port" default:"3306"`
		Username string `yaml:"username" default:"root"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname" default:"blocksbom"`
	} `yaml:"database"`

	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`

	Fabric struct {
		ConfigPath    string `yaml:"configPath"`
		ChannelID     string `yaml:"channelID"`
		ChaincodeName string `yaml:"chaincodeName"`
	} `yaml:"fabric"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// 获取环境变量中的配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认配置文件路径
		configPath = "../../configs/config.yaml"
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	if err := loadYamlConfig(configPath, cfg); err != nil {
		return nil, err
	}

	// 如果JWT密钥为空，生成新密钥并更新配置文件
	if cfg.JWT.Secret == "" {
		secret, err := utils.GenerateRandomSecret()
		if err != nil {
			return nil, fmt.Errorf("生成密钥失败: %v", err)
		}
		cfg.JWT.Secret = secret

		// 将更新后的配置写回文件
		newData, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("序列化配置失败: %v", err)
		}

		if err := os.WriteFile(configPath, newData, 0644); err != nil {
			return nil, fmt.Errorf("更新配置文件失败: %v", err)
		}
	}

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return cfg, nil
}

// validate 验证配置
func (c *Config) validate() error {
	if c.Database.Username == "" {
		return errors.New("数据库用户名不能为空")
	}
	if c.Database.Host == "" {
		return errors.New("数据库主机不能为空")
	}
	if c.Database.Port == 0 {
		c.Database.Port = 3306 // 使用默认端口
	}
	if c.Database.DBName == "" {
		return errors.New("数据库名不能为空")
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080 // 使用默认端口
	}
	return nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
	)
}

// loadYamlConfig 从YAML文件加载配置
func loadYamlConfig(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
