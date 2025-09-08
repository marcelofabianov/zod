package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	User     UserConfig     `mapstructure:"user"`
	General  GeneralConfig  `mapstructure:"general"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Server   ServerConfig   `mapstructure:"server"`
	DB       DBConfig       `mapstructure:"db"`
	Goose    GooseConfig    `mapstructure:"goose"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type UserConfig struct {
	UID int `mapstructure:"host_uid"`
	GID int `mapstructure:"host_gid"`
}

type GeneralConfig struct {
	Env string `mapstructure:"env"`
	TZ  string `mapstructure:"tz"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type ServerConfig struct {
	API  APIConfig  `mapstructure:"api"`
	CORS CORSConfig `mapstructure:"cors"`
}

type APIConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	RateLimit    int           `mapstructure:"rate_limit"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	MaxBodySize  int           `mapstructure:"maxbodysize"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedorigins"`
	AllowedMethods   []string `mapstructure:"allowedmethods"`
	AllowedHeaders   []string `mapstructure:"allowedheaders"`
	ExposedHeaders   []string `mapstructure:"exposedheaders"`
	AllowCredentials bool     `mapstructure:"allowcredentials"`
}

type DBConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"maxopenconns"`
	MaxIdleConns    int           `mapstructure:"maxidleconns"`
	ConnMaxLifetime time.Duration `mapstructure:"connmaxlifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connmaxidletime"`
	QueryTimeout    time.Duration `mapstructure:"querytimeout"`
	ExecTimeout     time.Duration `mapstructure:"exectimeout"`
}

type GooseConfig struct {
	Driver       string `mapstructure:"driver"`
	MigrationDir string `mapstructure:"migration_dir"`
	DBString     string `mapstructure:"dbstring"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	ReadTimeout  time.Duration `mapstructure:"readtimeout"`
	WriteTimeout time.Duration `mapstructure:"writetimeout"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetDefault("user.host_uid", 1000)
	v.SetDefault("user.host_gid", 1000)
	v.SetDefault("general.env", "development")
	v.SetDefault("general.tz", "UTC")
	v.SetDefault("logger.level", "info")
	v.SetDefault("server.api.host", "0.0.0.0")
	v.SetDefault("server.api.port", 8080)
	v.SetDefault("server.api.rate_limit", 100)
	v.SetDefault("server.api.read_timeout", "5s")
	v.SetDefault("server.api.write_timeout", "10s")
	v.SetDefault("server.api.idle_timeout", "120s")
	v.SetDefault("server.api.maxbodysize", 1048576)
	v.SetDefault("server.cors.allowedorigins", []string{"*"})
	v.SetDefault("server.cors.allowedmethods", []string{"GET", "POST"})
	v.SetDefault("server.cors.allowedheaders", []string{"Content-Type", "Authorization"})
	v.SetDefault("server.cors.exposedheaders", []string{})
	v.SetDefault("server.cors.allowcredentials", true)
	v.SetDefault("db.driver", "postgres")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", 5432)
	v.SetDefault("db.user", "user")
	v.SetDefault("db.password", "password")
	v.SetDefault("db.name", "app")
	v.SetDefault("db.ssl_mode", "disable")
	v.SetDefault("db.maxopenconns", 10)
	v.SetDefault("db.maxidleconns", 5)
	v.SetDefault("db.connmaxlifetime", "1h")
	v.SetDefault("db.connmaxidletime", "10m")
	v.SetDefault("db.querytimeout", "5s")
	v.SetDefault("db.exectimeout", "3s")
	v.SetDefault("goose.driver", "postgres")
	v.SetDefault("goose.migration_dir", "./db/migrations")
	v.SetDefault("goose.dbstring", "")
	v.SetDefault("rabbitmq.host", "localhost")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.user", "guest")
	v.SetDefault("rabbitmq.password", "guest")
	v.SetDefault("rabbitmq.vhost", "/")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.readtimeout", "3s")
	v.SetDefault("redis.writetimeout", "3s")

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
