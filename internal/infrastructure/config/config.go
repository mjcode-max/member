package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Log       LogConfig       `mapstructure:"log"`
	Auth      AuthConfig      `mapstructure:"auth"`
	Admin     AdminConfig     `mapstructure:"admin"`
	HuaweiFRS HuaweiFRSConfig `mapstructure:"huawei_frs"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type            string        `mapstructure:"type"`              // 数据库类型: mysql, postgres, file
	Host            string        `mapstructure:"host"`              // 数据库主机（mysql/postgres使用）
	Port            int           `mapstructure:"port"`              // 数据库端口（mysql/postgres使用）
	User            string        `mapstructure:"user"`              // 数据库用户（mysql/postgres使用）
	Password        string        `mapstructure:"password"`          // 数据库密码（mysql/postgres使用）
	DBName          string        `mapstructure:"dbname"`            // 数据库名称（mysql/postgres使用）
	Charset         string        `mapstructure:"charset"`           // 字符集（mysql使用）
	ParseTime       bool          `mapstructure:"parse_time"`        // 解析时间（mysql使用）
	Loc             string        `mapstructure:"loc"`               // 时区
	FilePath        string        `mapstructure:"file_path"`         // 文件路径（file类型使用）
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大生存时间
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret    string `mapstructure:"jwt_secret"`
	TokenExpires int    `mapstructure:"token_expires"` // 秒
}

// AdminConfig 默认管理员配置
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// HuaweiFRSConfig 华为云FRS配置
type HuaweiFRSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	ProjectID       string `mapstructure:"project_id"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	FaceSetName     string `mapstructure:"face_set_name"`
}

func (cfg *Config) GetJwtSecret() string {
	return cfg.Auth.JWTSecret
}

func (cfg *Config) GetTokenExpires() int {
	return cfg.Auth.TokenExpires
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 支持环境变量
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 转换时间单位为秒
	config.Server.ReadTimeout = config.Server.ReadTimeout * time.Second
	config.Server.WriteTimeout = config.Server.WriteTimeout * time.Second
	config.Database.ConnMaxLifetime = config.Database.ConnMaxLifetime * time.Second

	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	viper.SetDefault("database.type", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.file_path", "test.db")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	viper.SetDefault("auth.jwt_secret", "your-secret-key-change-in-production")
	viper.SetDefault("auth.token_expires", 7200) // 2小时

	viper.SetDefault("admin.username", "admin")
	viper.SetDefault("admin.password", "admin123")
}
