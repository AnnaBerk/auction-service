package app

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Env              string     `toml:"env"`
	ConnectionString string     `toml:"-"`
	Postgres         Postgres   `toml:"postgres"`
	HTTPServer       HTTPServer `toml:"http_server"`
	GRPCPort         string     `toml:"grpc_port" env:"GRPC_PORT" env-required:"true"`
}

type Postgres struct {
	Host        string        `toml:"host"`
	DBPort      int           `toml:"port"`
	User        string        `toml:"user"`
	Password    string        `toml:"password"`
	DBName      string        `toml:"dbname"`
	SSLMode     string        `toml:"sslmode"`
	Timeout     time.Duration `toml:"timeout"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
	MaxPoolSize int           `toml:"max_pool_size"`
}

type HTTPServer struct {
	Address     string        `toml:"address"`
	Timeout     time.Duration `toml:"timeout"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
	Port        string        `toml:"port"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	cfg.ConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.DBPort, cfg.Postgres.DBName, cfg.Postgres.SSLMode)

	return &cfg, nil
}
