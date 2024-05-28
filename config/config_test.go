package config_test

import (
	"testing"

	"github.com/lab42/mirror/config"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedConfig(t *testing.T) {
	// Read the embedded YAML file
	data, err := config.Default.ReadFile("config.yaml")
	assert.NoError(t, err)

	// Assert that the data is not empty
	assert.NotEmpty(t, data)
}
