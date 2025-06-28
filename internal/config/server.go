package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	devEnv  = "local"
	prodEnv = "prod"
	testEnv = "test"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	HTTP     HTTPConfig     `yaml:"http"`
	GH       GHConfig       `yaml:"github"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type AppConfig struct {
	// TODO: check wether cleanenv can parse custom type like envType for ENV var
	ENV                 string `yaml:"env" env:"APP_ENV" env-required:"true"`
	GithubWebhookSecret string `yaml:"github_app_secret" env:"WEBHOOK_SECRET" env-required:"true"`
}

type HTTPConfig struct {
	Port           int           `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	ReadTimeout    time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
	WriteTimeout   time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
	GatewayTimeout time.Duration `yaml:"gateway_timeout" env:"HTTP_GATEWAY_TIMEOUT" env-required:"true"`
}

type RabbitMQConfig struct {
	Host               string `yaml:"host" env:"RABBITMQ_HOST" env-required:"true"`
	Port               int    `yaml:"port" env:"RABBITMQ_PORT" env-required:"true"`
	User               string `yaml:"user" env:"RABBITMQ_USER" env-required:"true"`
	Pass               string `yaml:"pass" env:"RABBITMQ_PASS" env-required:"true"`
	IssueExchange      string `yaml:"issue_exchange" env:"RABBITMQ_ISSUE_EXCHANGE" env-default:"issue_exchange"`
	CheckRequestsQueue string `yaml:"check_requests_queue" env:"RABBITMQ_CHECK_REQUESTS_QUEUE" env-default:"manual_check_requests"` // TODO: remove
}

func (c *RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.User, c.Pass, c.Host, c.Port)
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
