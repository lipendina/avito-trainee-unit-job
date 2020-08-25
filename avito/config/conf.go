package config

import (
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const confPath = "config/parameters.yaml"

type DBConfig struct {
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
	Host     string `yaml:"db_host"`
	Port     uint16 `yaml:"db_port"`
	DBName   string `yaml:"db_name"`
}

type ApplicationConfig struct {
	DB DBConfig `yaml:",inline"`
	HTTPPort uint16 `yaml:"http_port"`
}

func ParseConfig() (*ApplicationConfig, error) {
	config := &ApplicationConfig{}

	confFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, xerrors.Errorf("Failed to read config file: %+v", err)
	}

	confFile = []byte(os.ExpandEnv(string(confFile)))

	if err = yaml.Unmarshal(confFile, config); err != nil {
		return nil, xerrors.Errorf("Cannot unmarshal config: %+v", err)
	}

	return config, nil
}
