package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
)

var k = koanf.New(".")

func Load(configPath string) (*Config, error) {
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, errors.Wrap(err, "unable to load config")
	}
	config := &Config{}
	if err := k.Unmarshal("", config); err != nil {
		return nil, errors.Wrap(err, "unable to load config")
	}

	return config, nil
}

const (
	LOAD_BALANCING_ROUND_ROBIN = "LOAD_BALANCING_ROUND_ROBIN"
)

type Config struct {
	Server        *Server
	RoutesMapping map[string]*TargetRoutes
}

type Server struct {
	Port int
}

type TargetRoutes struct {
	Origins               []string
	LoadBalancingStrategy string
}
