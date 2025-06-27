package config

import "github.com/ilyakaznacheev/cleanenv"

type Client struct {
}

func MustLoadClient() *Client {
	var client Client
	if err := cleanenv.ReadEnv(&client); err != nil {
		panic("failed to load environment variables: " + err.Error())
	}

	return &client
}
