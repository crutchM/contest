package domain

type Config struct {
	Port    int    `env:"PORT"`
	ApiKey  string `env:"API_KEY"`
	ConnStr string `env:"CONN_STR"`
}
