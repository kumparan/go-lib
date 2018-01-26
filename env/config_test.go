package env_test

import (
	"testing"

	"github.com/lab46/example/pkg/env"
)

func TestSetAndGetConfigDir(t *testing.T) {
	dir := "../../files/config/testconfig"
	err := env.SetConfigDir(dir)
	if err != nil {
		t.Error(err)
	}
	confdir := env.GetConfigDir()
	if dir != confdir {
		t.Errorf("Expecting %s but got %s", dir, confdir)
	}
}

func TestLoadYamlConfig(t *testing.T) {
	configDir := "../../files/config/testconfig"
	err := env.SetConfigDir(configDir)
	if err != nil {
		t.Error(err)
	}
	tc := struct {
		test struct {
			key1 string `yaml:"key1"`
		}
	}{}
	err = env.LoadYamlConfig(&tc, "test.yml")
	if err != nil {
		t.Error(err)
	}
}
