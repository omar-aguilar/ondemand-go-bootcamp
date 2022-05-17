package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	baseConfig := GetConfig()
	expectedDBPath := "./data/db.csv"
	assert.Equal(t, expectedDBPath, baseConfig.GetDBPath())
}
