package env_test

import (
	"testing"

	"github.com/lab46/example/pkg/env"
)

func TestSetFromEnvFile(t *testing.T) {
	err := env.SetFromEnvFile("../../files/.env")
	if err != nil {
		t.Error(err)
	}
	val1 := env.Getenv("KEY1")
	if val1 != "value1" {
		t.Errorf("Invalid KEY1 value. Value: %s", val1)
	}
	val2 := env.Getenv("KEY2")
	if val2 != "value2" {
		t.Errorf("Invalid KEY2 value. Value: %s", val2)
	}
	val3 := env.Getenv("KEY3")
	if val3 != "value3" {
		t.Errorf("Invalid KEY3 value. Value: %s", val3)
	}
}

func TestGetCurrentServieEnv(t *testing.T) {
	e := env.GetCurrentServiceEnv()
	if e != "dev" {
		t.Errorf("Current env should be development. Current env: %s", e)
	}

	err := env.Setenv("EXMPLENV", "prod")
	if err != nil {
		t.Error(err)
	}
	e = env.GetCurrentServiceEnv()
	if e != "prod" {
		t.Errorf("Current env should changed to prod. Current env: %s", e)
	}
}

func TestGetSet(t *testing.T) {
	cases := []struct {
		key   string
		value string
	}{
		{
			key:   "key1",
			value: "value1",
		},
		{
			key:   "key2",
			value: "value2",
		},
		{
			key:   "key3",
			value: "value3",
		},
	}

	for _, c := range cases {
		err := env.Setenv(c.key, c.value)
		if err != nil {
			t.Error(err)
		}
		val := env.Getenv(c.key)
		if val != c.value {
			t.Errorf("Expecting %s but got %s", c.value, val)
		}
	}
}
