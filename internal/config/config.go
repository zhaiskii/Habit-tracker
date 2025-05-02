package config

import(
	"github.com/ilyakaznacheev/cleanenv"
)

type Http struct{
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct{
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Name string `yaml:"name"`
	Pass string `yaml:"password"`
}

type Config struct {
	Http `yaml:"server"`
	Database `yaml:"database"`
}

func MustLoad() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("./idkhowtonamethat.yaml", &cfg)
	if err!=nil{
		return nil, err
	}
	return &cfg, nil
}