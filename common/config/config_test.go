package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const appName = "testapp"

//nolint:errcheck
func TestLoadConfigFromEnvironment(t *testing.T) {
	// setup
	os.Setenv("TESTAPP_FOO", "foo_value")
	os.Setenv("TESTAPP_BAR", "bar_value")
	defer func() {
		os.Unsetenv("FOO")
		os.Unsetenv("BAR")
	}()

	v, _ := LoadConfig(appName)

	assert.Equal(t, "foo_value", v.GetString("FOO"))
	assert.Equal(t, "bar_value", v.GetString("bar"))
}

//nolint:errcheck
func TestLoadConfigFromFile(t *testing.T) {
	// setup
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "testapp.yaml")
	configContent := []byte(`
foo: "foo_value"
bar: "bar_value"
`)
	err := os.WriteFile(configFile, configContent, 0644)
	require.NoError(t, err)

	// create symlink to make the config accessible in the current directory
	err = os.Symlink(configFile, "./testapp.yaml")
	require.NoError(t, err)
	defer os.Remove("./testapp.yaml")

	v, err := LoadConfig(appName)

	require.NoError(t, err)
	assert.Equal(t, "foo_value", v.GetString("foo"))
	assert.Equal(t, "bar_value", v.GetString("bar"))
}

//nolint:errcheck
func TestLoadConfigEnvPriorityOverFile(t *testing.T) {
	// create the config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "testapp.yaml")
	configContent := []byte(`
foo: "foo_value_file"
bar: "bar_value_file"
`)
	err := os.WriteFile(configFile, configContent, 0644)
	require.NoError(t, err)

	// Create symlink to make the config accessible in current directory
	err = os.Symlink(configFile, "./testapp.yaml")
	require.NoError(t, err)
	defer os.Remove("./testapp.yaml")

	// Set environment variables
	os.Setenv("TESTAPP_FOO", "foo_value_file")
	os.Setenv("TESTAPP_BAR", "bar_value_env")
	defer func() {
		os.Unsetenv("TESTAPP_FOO")
		os.Unsetenv("TESTAPP_BAR")
	}()

	v, err := LoadConfig(appName)

	require.NoError(t, err)
	assert.Equal(t, "foo_value_file", v.GetString("foo"))
	assert.Equal(t, "bar_value_env", v.GetString("bar"))
}
