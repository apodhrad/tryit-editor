package service

import (
	"bytes"
	"errors"
	"os/exec"
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
