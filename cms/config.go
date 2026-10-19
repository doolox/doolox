package cms

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config is the top-level configuration file (~/doolox/config.json). It only
// names the current website; per-website settings live in each website's own
// config.json (see SiteConfig).
type Config struct {
	Website string `json:"website"`
}

// SiteConfig is a per-website configuration file
// (~/doolox/websites/<site>/config.json).
type SiteConfig struct {
	Theme string `json:"theme"`
	URL   string `json:"url"`
}

// loadJSON decodes the JSON file at path into v.
func loadJSON(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[config] open %s: %v", path, err)
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(v); err != nil {
		log.Printf("[config] decode %s: %v", path, err)
		return err
	}
	return nil
}

func initConfig() *Config {
	c := &Config{}
	if err := loadJSON(filepath.Join(dataRoot(), "config.json"), c); err != nil {
		log.Printf("[initConfig] could not load config: %v", err)
	}
	return c
}

// siteConfig loads the per-website configuration for the named website. It
// returns an empty SiteConfig when the file is missing or unreadable.
func siteConfig(site string) *SiteConfig {
	c := &SiteConfig{}
	if err := loadJSON(filepath.Join(websiteDir(site), "config.json"), c); err != nil {
		log.Printf("[siteConfig] could not load config for %q: %v", site, err)
	}
	return c
}

// defaultWebsite is the website name used when config.json does not set one.
const defaultWebsite = "simple"

// Website returns the name of the current website from the configuration. It
// falls back to defaultWebsite when no name is configured.
func Website() string {
	if conf != nil && conf.Website != "" {
		return conf.Website
	}
	return defaultWebsite
}
