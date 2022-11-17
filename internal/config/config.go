package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var cfg config

type config struct {
	ServerAddr string `yaml:"server_addr"`
	PSQLLogin  string `yaml:"psql_login"`
	PSQLPass   string `yaml:"psql_password"`
	PSQLPort   string `yaml:"psql_port"`
	DBHost     string `yaml:"db_host"`
}

func InitConfig(args []string) error {

	var configPath string

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	flags.StringVar(&configPath, "c", "/opt/config.yaml", "set path to config")
	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}

	clean := filepath.Clean(configPath)

	file, err := os.Open(clean)
	if err != nil {
		return fmt.Errorf("fail to open config file in path \"%s\" with error %w", configPath, err)
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return fmt.Errorf("fail to parse config %w", err)
	}

	return nil
}

func GetConfig() config {
	return cfg
}
