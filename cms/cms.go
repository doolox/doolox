package cms

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *Config

var db *gorm.DB

func init() {
	setupDataDir()

	conf = initConfig()

	db = initDb()

	m = initMacaron()
}

// Start runs the Doolox web server and shows the system tray icon. It blocks
// until the tray is quit, then shuts the server down and returns.
func Start(trayIcon TrayIcon) {
	m.Get("/", pageView)
	m.Get("/generate", generateView)
	m.Get("/:page", pageView)

	server := &http.Server{
		Addr:    "0.0.0.0:5000",
		Handler: m,
	}

	go func() {
		log.Printf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[Start] server error: %v", err)
		}
	}()

	runTray(trayIcon, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[Start] server shutdown error: %v", err)
		}
	})
}
