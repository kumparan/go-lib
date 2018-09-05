package env

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kumparan/go-lib/logger"
	"os"
	"runtime"
	"strings"
)

type ServiceEnv string

const (
	DevelopmentEnv ServiceEnv = "dev"
	StagingEnv     ServiceEnv = "staging"
	ProductionEnv  ServiceEnv = "prod"
)

var (
	envName      = "EXMPLENV"
	currentBuild = "unavailable"
	goVersion    string
)

// env package will read .env file when applicatino is started

func init() {
	err := SetFromEnvFile(".env")
	if err != nil {
		logger.Debug(err)
	}
	goVersion = runtime.Version()
}

func SetFromEnvFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	if err := scanner.Err(); err != nil {
		return err
	}
	for scanner.Scan() {
		text := scanner.Text()
		vars := strings.SplitN(text, "=", 2)
		if len(vars) < 2 {
			return err
		}
		if err := Setenv(vars[0], vars[1]); err != nil {
			return err
		}
	}
	return nil
}

func SetEnvName(name string) {
	envName = name
}

func EnvList() []ServiceEnv {
	return []ServiceEnv{DevelopmentEnv, StagingEnv, ProductionEnv}
}

func SetCurrentServiceEnv(env ServiceEnv) error {
	return Setenv(envName, string(env))
}

func GetCurrentServiceEnv() string {
	e := Getenv(envName)
	if e == "" {
		e = string(DevelopmentEnv)
	}
	return e
}

func Getenv(key string) string {
	return os.Getenv(key)
}

func Setenv(key, value string) error {
	return os.Setenv(key, value)
}

// SetCurrentBuild to determine the latest build of
func SetCurrentBuild(buildnumber string) {
	currentBuild = buildnumber
}

// GetCurrentBuild return the current build number
func GetCurrentBuild() string {
	return currentBuild
}

// GetGoVersion to return current build go version
func GetGoVersion() string {
	return goVersion
}

// LoadEnv Load environment from a file name set in in an environment
func LoadEnv(envName string) (err error) {
	envFile, envExist := os.LookupEnv(envName)
	if envExist {
		fmt.Printf("Loading environment from %s\n", envFile)
		err = godotenv.Load(envFile)
	}
}
