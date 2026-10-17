package cms

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
// config.json lives in dataRoot.
func websitesDir() string { return filepath.Join(dataRoot(), "websites") }

// websites returns the names of every available website: the default website
// plus each of its subfolders that contains a site source layout (either the
// modern src/ directory or the legacy content/ directory). The result is
// sorted.
func websites() []string {
	// The default website is always available, served at the root.
	names := []string{Website()}

	main := filepath.Join(websitesDir(), Website())
	entries, err := os.ReadDir(main)
	if err != nil {
		return names
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		// Nested websites are subfolders of the default website that contain
		// their own site source layout.
		dir := filepath.Join(main, entry.Name())
		if dirExists(filepath.Join(dir, "src")) || dirExists(filepath.Join(dir, "content")) {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

// isWebsite reports whether name is one of the available websites.
func isWebsite(name string) bool {
	for _, n := range websites() {
		if n == name {
			return true
		}
	}
	return false
}

// dirExists reports whether path is an existing directory.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// websiteDir is the folder of the named website. The default website lives
// directly under the websites directory; every other website is nested inside
// the default website's folder.
func websiteDir(name string) string {
	if name == Website() {
		return filepath.Join(websitesDir(), name)
	}
	return filepath.Join(websitesDir(), Website(), name)
}

// sourceDir is the current source layout for a website.
func sourceDir(name string) string { return filepath.Join(websiteDir(name), "src") }

// templatesDir is a folder reserved for templates (currently unused).
func templatesDir() string { return filepath.Join(dataRoot(), "templates") }

// contentDir is the directory macaron loads templates from for the named
// website. The legacy content layout remains supported for existing
// installations.
func contentDir(name string) string {
	if dirExists(sourceDir(name)) {
		return sourceDir(name)
	}
	return filepath.Join(websiteDir(name), "content")
}

// assetsDir returns the static asset directory for the named website.
func assetsDir(name string) string {
	if contentDir(name) == sourceDir(name) {
		return filepath.Join(sourceDir(name), "assets")
	}
	return filepath.Join(contentDir(name), "theme", "assets")
}

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

// setupWebsiteDir creates the root folder for the current website.
// It must run after initDb, because the website name is read from the
// key-value store.
func setupWebsiteDir() {
	if err := os.MkdirAll(websiteDir(Website()), 0755); err != nil {
		log.Printf("[setup] could not create %s: %v", websiteDir(Website()), err)
	}
}
