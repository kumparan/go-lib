package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/kumparan/go-lib/log"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Dir string
	Env string
}

// defaut value of config
var cfg = config{
	Dir: "../../files/config",
	Env: string(DevelopmentEnv),
}

// SetConfigDir to set a config directory
func SetConfigDir(dir string) error {
	if f, err := os.Stat(dir); err != nil {
		return fmt.Errorf("failed to check directory: %s", err.Error())
	} else if !f.IsDir() {
		return fmt.Errorf("%s is not a dir: %s", dir, err.Error())
	}
	cfg.Dir = dir
	return nil
}

// GetConfigDir return current config directory
func GetConfigDir() string {
	return cfg.Dir
}

// LoadYamlConfig to load config from a yaml file format
func LoadYamlConfig(result interface{}, filename string) error {
	e := GetCurrentServiceEnv()
	dirEnv := strings.ToLower(e)
	if dirEnv == "" {
		dirEnv = string(DevelopmentEnv)
	}
	confDir := path.Join(cfg.Dir, dirEnv, filename)
	log.Debugf("[config][yaml] from: %s", confDir)
	content, err := ioutil.ReadFile(confDir)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, result)
}
