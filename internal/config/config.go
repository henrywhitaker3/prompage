package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrUnknownDatasource = errors.New("unknown datasource for query")
)

type Query struct {
	Name       string        `yaml:"name"`
	Query      string        `yaml:"query"`
	Expression string        `yaml:"expression"`
	Range      time.Duration `yaml:"range"`
	Step       time.Duration `yaml:"step"`
	BoolValue  bool          `yaml:"bool"`
	Units      string        `yaml:"units"`
	Datasource string        `yaml:"datasource"`
}

type Service struct {
	Name   string  `yaml:"name"`
	Query  Query   `yaml:"query"`
	Group  string  `yaml:"group"`
	Extras []Query `yaml:"extras"`
}

type UI struct {
	PageTitle       string        `yaml:"title"`
	RefreshInterval time.Duration `yaml:"refresh"`
	Graphs          Graphs        `yaml:"graphs"`
}

type Metrics struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type Graphs struct {
	Points int `yaml:"points"`
}

type Datasource struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Url    string            `yaml:"url"`
	Extras map[string]string `yaml:"extras"`
}

type Config struct {
	Port     int       `yaml:"port"`
	Metrics  Metrics   `yaml:"metrics"`
	Services []Service `yaml:"services"`

	Datasources []Datasource `yaml:"datasources"`

	Prometheus string `yaml:"prometheus"`

	Refresh time.Duration `yaml:"refresh"`

	UI UI `yaml:"ui"`
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

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

func setDefaults(conf *Config) {
	if conf.Port == 0 {
		conf.Port = 3000
	}
	if conf.Metrics.Port == 0 {
		conf.Metrics.Port = 9743
	}

	if conf.Prometheus != "" {
		conf.Datasources = append(conf.Datasources, Datasource{
			Name: "prometheus",
			Type: "prometheus",
			Url:  conf.Prometheus,
		})
	}

	for i, svc := range conf.Services {
		svc.Query.Name = "main"
		if svc.Group == "" {
			svc.Group = "default"
		}
		svc.Query.BoolValue = true
		setDefaultQueryValues(&svc.Query)

		for i, query := range svc.Extras {
			if query.Expression == "" {
				query.Expression = "float(result)"
			}
			setDefaultQueryValues(&query)
			svc.Extras[i] = query
		}
		conf.Services[i] = svc
	}

	if conf.Refresh == 0 {
		conf.Refresh = time.Second * 30
	}

	if conf.UI.PageTitle == "" {
		conf.UI.PageTitle = "PromPage"
	}
	if conf.UI.RefreshInterval == 0 {
		conf.UI.RefreshInterval = time.Second * 30
	}
	if conf.UI.Graphs.Points <= 0 {
		conf.UI.Graphs.Points = 200
	}
}

func setDefaultQueryValues(q *Query) {
	if q.Range == 0 {
		q.Range = time.Hour * 24
	}
	if q.Step == 0 {
		q.Step = time.Minute * 5
	}
	if q.Datasource == "" {
		q.Datasource = "prometheus"
	}
}

func (c *Config) Validate() error {
	if c.Prometheus == "" && len(c.Datasources) == 0 {
		return errors.New("you must configure a datasource")
	}

	if c.Prometheus != "" {
		log.Println("DEPRECATED - the prometheus config option is deprecated, replace with an entry in datasources with name prometheus")
	}

	for _, ds := range c.Datasources {
		if ds.Name == "" {
			return errors.New("all datasources must have a name")
		}
		if !slices.Contains([]string{"prometheus", "datadog"}, ds.Type) {
			return errors.New("datasources must be one of: prometheus, datadog")
		}
		if ds.Url == "" {
			return errors.New("datasources must have a url configured")
		}
		if ds.Type == "datadog" {
			if _, ok := ds.Extras["apiKey"]; !ok {
				return errors.New("datadog extra apiKey must be set")
			}
			if _, ok := ds.Extras["appKey"]; !ok {
				return errors.New("datadog extra appKey must be set")
			}
		}
	}

	for _, svc := range c.Services {
		if !c.containsDatasource(svc.Query.Datasource) {
			return ErrUnknownDatasource
		}
		for _, extra := range svc.Extras {
			if extra.Name == "" {
				return errors.New("extra query name cannot be empty")
			}
			if !c.containsDatasource(extra.Datasource) {
				return ErrUnknownDatasource
			}
		}
	}

	return nil
}

func (c *Config) containsDatasource(name string) bool {
	for _, ds := range c.Datasources {
		if name == ds.Name {
			return true
		}
	}
	return false
}
