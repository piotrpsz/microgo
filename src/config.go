package src

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

// Configuration values used acros the app
type Configuration struct {
	Gin struct {
		ReleaseMode bool `json:"release" env:"GIN_RELEASE" env-default:"false"`
	}
	Logger struct {
		Enabled     bool   `json:"enabled" env:"LOG_ENABLED" env-required:"false"`
		FileName    string `json:"file_name" env:"LOG_FILE_NAME" env-default:"mend.cfg"`
		DirName     string `json:"dir_name" env:"LOG_DIR_NAME"`
		RootDirName string `json:"root_dir" env:"LOG_ROOT_DIR"`
	}
}

var (
	instance *Configuration
	once     sync.Once
)

// Config returns instance of Configuration structure.
// Configuration is implemented as singleton!
func Config() *Configuration {
	once.Do(func() {
		instance = getcfg()
	})
	return instance
}

// getcfg reads and sets config data
// the highest priority has value from .env file
// next environment variable
// or default value (env-default in structure's tags)
func getcfg() *Configuration {
	cfg, err := load()
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// load fetches config values
func load() (*Configuration, error) {
	cfg := new(Configuration)

	// firstly try to read file '.env'
	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if err != nil {
		// when something was wrong read
		// environment variables
		err = cleanenv.ReadEnv(cfg)
	}
	return cfg, err
}

// String conform to stringify
func (c *Configuration) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Error(err)
		return "?"
	}
	return string(data)
}

// LogFilePath returns comple path to log file.
// When is needed directory components are created.
// WARNING: log file is for only one session (old file id deleted).
func (c *Configuration) LogFilePath() (string, error) {
	if c.Logger.RootDirName == "" {
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		c.Logger.RootDirName = dir
	}

	fpath := c.Logger.RootDirName
	if c.Logger.DirName != "" {
		fpath = filepath.Join(fpath, c.Logger.DirName)
		if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
			return "", err
		}
	}

	fpath = filepath.Join(fpath, c.Logger.FileName)

	// delete previous instance of the log file
	// (if exists)
	_ = os.Remove(fpath)

	return fpath, nil
}
