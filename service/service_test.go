package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCat(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "cat-*-input")
	assert.Nil(t, err)
	_, err = tmpFile.WriteString("Hello Test")
	assert.Nil(t, err)
	err = tmpFile.Close()
	assert.Nil(t, err)

	svc := SERVICE_CAT
	assert.Equal(t, "cat", svc.Name())
	out, err := svc.Run(tmpFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, "Hello Test", string(out))

	svc = Service{Exec: "non-existing-svc"}
	assert.Equal(t, "non-existing-svc", svc.Name())
	out, err = svc.Run(tmpFile.Name())
	assert.NotNil(t, err)
	assert.Equal(t, "exec: \"non-existing-svc\": executable file not found in $PATH", err.Error())
	assert.Equal(t, "", string(out))
}
