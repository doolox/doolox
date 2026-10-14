package cms

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// dooloxHomeEnv overrides the base data directory (mainly useful for tests).
const dooloxHomeEnv = "DOOLOX_HOME"

// dataRoot returns the base directory holding all of Doolox's writable data.
// It defaults to ~/doolox and can be overridden with the DOOLOX_HOME
// environment variable.
func dataRoot() string {
	if dir := os.Getenv(dooloxHomeEnv); dir != "" {
		return dir
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		home = "."
	}
	return filepath.Join(home, "doolox")
}

// websitesDir is where the site data lives: content, docs, config.json and
// doolox.db.
func websitesDir() string { return filepath.Join(dataRoot(), "websites") }

// templatesDir is a folder reserved for templates (currently unused).
func templatesDir() string { return filepath.Join(dataRoot(), "templates") }

// contentDir is the directory macaron loads templates from.
func contentDir() string { return filepath.Join(websitesDir(), "content") }

// docsDir is the folder the static website is generated into.
func docsDir() string { return filepath.Join(websitesDir(), "docs") }

// logFile is the file all log output is written to.
func logFile() string { return filepath.Join(dataRoot(), "doolox.log") }

// setupDataDir creates the data directory layout and redirects logging to
// ~/doolox/doolox.log (mirrored to stderr). It must run before any config or
// database file is opened. The content directory is intentionally not created
// here: it is seeded by the installer/launcher on first run.
func setupDataDir() {
	for _, dir := range []string{websitesDir(), docsDir(), templatesDir()} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("[setup] could not create %s: %v", dir, err)
		}
	}

	f, err := os.OpenFile(logFile(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[setup] could not open log file %s: %v", logFile(), err)
		return
	}
	log.SetOutput(io.MultiWriter(f, os.Stderr))
}
