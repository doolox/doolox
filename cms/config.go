package cms

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config struct holds all our configuration
type Config struct {
	Theme   string `json:"theme"`
	URL     string `json:"url"`
	Website string `json:"website"`
}

// Load method loads configuration file to Config struct
func (sc *Config) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Got error while decoding JSON: %v", err)
		return err
	}

	return nil
}

func initConfig() *Config {
	c := &Config{}
	if err := c.Load(filepath.Join(dataRoot(), "config.json")); err != nil {
		log.Printf("[initConfig] could not load config: %v", err)
	}
	return c
}

// defaultWebsite is the website name used when config.json does not set one.
const defaultWebsite = "doolox"

// Website returns the name of the current website from the configuration. It
// falls back to defaultWebsite when no name is configured.
func Website() string {
	if conf != nil && conf.Website != "" {
		return conf.Website
	}
	return defaultWebsite
}
