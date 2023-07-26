package service

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"
)

var SERVICE_CAT Service = Service{
	Exec:     "cat",
	Examples: []string{"examples/default.txt"},
}

type Service struct {
	Exec     string
	Examples []string
}

func (s Service) Name() string {
	return s.Exec
}

func (s Service) Run(inputFile string) ([]byte, error) {
	var outb, errb bytes.Buffer

	cmd := exec.Command(s.Exec, inputFile)
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		if errb.String() != "" {
			err = errors.New(errb.String() + "\n" + err.Error())
		}
		return []byte{}, err
	}
	return outb.Bytes(), nil
}

func LoadServices(configFile string) ([]Service, error) {
	var services []Service
	if configFile == "" {
		return services, errors.New("No config file")
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return services, err
	}
	err = yaml.Unmarshal(data, &services)
	if err != nil {
		return services, err
	}
	return services, nil
}
