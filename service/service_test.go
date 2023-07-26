package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tmpFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "tmp-*")
	assert.Nil(t, err)
	_, err = tmpFile.WriteString(content)
	assert.Nil(t, err)
	err = tmpFile.Close()
	assert.Nil(t, err)
	return tmpFile.Name()
}

func TestService(t *testing.T) {
	tmpFile := tmpFile(t, "Hello Test")

	svc := SERVICE_CAT
	assert.Equal(t, "cat", svc.Name())
	out, err := svc.Run(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, "Hello Test", string(out))

	svc = Service{Exec: "non-existing-svc"}
	assert.Equal(t, "non-existing-svc", svc.Name())
	out, err = svc.Run(tmpFile)
	assert.NotNil(t, err)
	assert.Equal(t, "exec: \"non-existing-svc\": executable file not found in $PATH", err.Error())
	assert.Equal(t, "", string(out))
}

func TestLoadingServices(t *testing.T) {
	tmpFile := tmpFile(t, "- exec: foo1\n  examples:\n  - example11\n- exec: foo2\n  examples:\n  - example21\n  - example22")
	services, err := LoadServices(tmpFile)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(services))
	assert.Equal(t, "foo1", services[0].Name())
	assert.Equal(t, "example11", services[0].Examples[0])
	assert.Equal(t, "foo2", services[1].Name())
	assert.Equal(t, "example21", services[1].Examples[0])
	assert.Equal(t, "example22", services[1].Examples[1])
}
