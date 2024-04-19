package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Query struct {
	Name  string `yaml:"name"`
	Query string `yaml:"query"`
}

type Service struct {
	Name        string  `yaml:"name"`
	Queries     []Query `yaml:"queries"`
	ShowMetrics bool    `yaml:"show_metrics"`
}

type Config struct {
	Port       int       `yaml:"port"`
	Services   []Service `yaml:"services"`
	Prometheus string    `yaml:"prometheus"`
}

func Load(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file stat: %v", err)
	}

	fb, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config file read: %v", err)
	}

	conf := &Config{}
	if err := yaml.Unmarshal(fb, conf); err != nil {
		return nil, fmt.Errorf("config file parse: %v", err)
	}

	setDefaults(conf)

	return conf, nil
}

func setDefaults(conf *Config) {
	if conf.Port == 0 {
		conf.Port = 3000
	}
}
