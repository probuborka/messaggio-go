package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigCorrect(t *testing.T) {
	cfgDir := "configs"
	cfgFile := "config_debug"
	cfg, err := Init(cfgDir, cfgFile)

	require.Empty(t, err)
	require.NotEmpty(t, cfg)

	//HTTP.Port
	assert.NotEmpty(t, cfg.HTTP.Port, fmt.Errorf("PORT is Empty"))
}

func TestDirectoryConfigNotFound(t *testing.T) {
	cfgDir := "configsNotFound"
	cfgFile := "config_debug"
	_, err := Init(cfgDir, cfgFile)

	require.NotEmpty(t, err)
	require.Equal(t, err.Error(), "directory not found")
}
