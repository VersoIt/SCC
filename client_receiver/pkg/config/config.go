package config

import (
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type Config struct {
	Host         string `yaml:"host"`
	WorkersCount int    `yaml:"workers_count"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		file, err := os.Open("config.yml")
		if err != nil {
			logrus.Errorf("error opening file: %s", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logrus.Errorf("error closing file: %s", err)
			}
		}(file)

		data, err := io.ReadAll(file)
		if err != nil {
			logrus.Errorf("error reading file: %s", err)
		}

		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			logrus.Errorf("error unmarshalling file: %s", err)
		}
	})
	return cfg
}
