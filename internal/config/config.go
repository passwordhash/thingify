package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App AppConfig `yaml:"app"`
	GH  GHConfig  `yaml:"github"`
}

type AppConfig struct {
	GHQueriesPath   string        `env:"GH_QUERIES_PATH" yaml:"github_queries_path" env-default:"./queries/github"`
	PollingInterval time.Duration `env:"POLLING_INTERVAL" yaml:"polling_interval" env-required:"true"`
}

type GHConfig struct {
	BaseURL string `yaml:"api_url" env:"GH_API_BASE_URL" env-default:"https://api.github.com"`
	Token   string `env:"GH_TOKEN" env-required:"true"`
}

// MustLoad загружает конфигурацию из файла и переменных окружения.
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is not set")
	}

	return MustLoadByPath(path)
}

// MustLoadByPath загружает конфигурацию из указанного файла.
// Если файл не существует или нет прав доступа, вызывает панику.
func MustLoadByPath(configPath string) *Config {
	_, err := os.Stat(configPath)
	if err != nil && os.IsPermission(err) {
		panic("no permission to config file: " + configPath)
	}
	if err != nil && os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to load config: " + err.Error())
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath получает путь к конфигурационному файлу из флага `config` или
// переменной окружения `CONFIG_PATH`. Если путь не указан, возвращает пустую строку.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
