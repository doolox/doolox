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

// websitesDir is where the websites live, one subfolder per website.
// config.json and doolox.db live in dataRoot.
func websitesDir() string { return filepath.Join(dataRoot(), "websites") }

// websiteDir is the folder of the current website, named after the website
// name stored in the key-value store. It holds the site content and the
// generated static site.
func websiteDir() string { return filepath.Join(websitesDir(), Website()) }

// templatesDir is a folder reserved for templates (currently unused).
func templatesDir() string { return filepath.Join(dataRoot(), "templates") }

// contentDir is the directory macaron loads templates from.
func contentDir() string { return filepath.Join(websiteDir(), "content") }

// docsDir is the folder the static website is generated into.
func docsDir() string { return filepath.Join(websiteDir(), "docs") }

// logFile is the file all log output is written to.
func logFile() string { return filepath.Join(dataRoot(), "doolox.log") }

// setupDataDir creates the base data directory layout and redirects logging to
// ~/doolox/doolox.log (mirrored to stderr). It must run before any config or
// database file is opened. The website content directory is intentionally not
// created here: it is seeded by the installer/launcher on first run.
func setupDataDir() {
	for _, dir := range []string{dataRoot(), websitesDir(), templatesDir()} {
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

// setupWebsiteDir creates the generated-site folder for the current website.
// It must run after initDb, because the website name is read from the
// key-value store.
func setupWebsiteDir() {
	if err := os.MkdirAll(docsDir(), 0755); err != nil {
		log.Printf("[setup] could not create %s: %v", docsDir(), err)
	}
}
