package usersdb

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

var connStr string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error getting the current file path.")
	}

	dir := filepath.Dir(filename)
	projectRoot := filepath.Join(dir, "..", "..", "..")
	configPath := filepath.Join(projectRoot, "config", "test.yml")

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading testing configuration file: " + err.Error())
	}

	connStr = viper.GetString("database.connection_string")
	if connStr == "" {
		panic("Connection string not found in configuration")
	}
}

func TestNewDB(t *testing.T) {

	// Test with txdb enabled
	t.Run("WithTxDB", func(t *testing.T) {
		_, err := NewDB(connStr, true, "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	// Test with txdb disabled
	t.Run("WithoutTxDB", func(t *testing.T) {
		_, err := NewDB(connStr, false, "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}
