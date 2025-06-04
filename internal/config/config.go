package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GH              GHConfig      `yaml:"github"`
	GHQueriesPath   string        `env:"GH_QUERIES_PATH" yaml:"github_queries_path" env-default:"./queries/github"`
	PollingInterval time.Duration `env:"POLLING_INTERVAL" yaml:"polling_interval" env-required:"true"`
}

type GHConfig struct {
	BaseURL string `yaml:"api_url" env:"GH_API_BASE_URL" env-default:"https://api.github.com"`
	Token   string `env:"GH_TOKEN" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config

	//if err := cleanenv.ReadConfig(, &cfg); err != nil {
	//	panic("failed to load config: " + err.Error())
	//}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}

	return &cfg
}
