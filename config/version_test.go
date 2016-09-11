package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	version := GetVersion()

	assert.Equal(t, "dev (HEAD)", version.ShortInfo())
}
