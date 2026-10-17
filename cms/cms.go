package cms

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *Config

func init() {
	setupDataDir()

	conf = initConfig()

	setupWebsiteDir()

	m = initMacaron()
}

// Start runs the Doolox web server and shows the system tray icon. It blocks
// until the tray is quit, then shuts the server down and returns.
func Start(trayIcon TrayIcon, openWebsite bool) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m.Get("/", pageView)
	m.Get("/generate", generateView)
	m.Get("/doolox.js", dooloxJSView)
	m.Get("/:page", pageView)

	server := &http.Server{
		Addr:    "127.0.0.1:5000",
		Handler: m,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Printf("[Start] server listen error: %v", err)
	} else {
		go func() {
			log.Printf("listening on %s", server.Addr)
			if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
				log.Printf("[Start] server error: %v", err)
			}
		}()

		if openWebsite {
			if err := openBrowser(websiteURL); err != nil {
				log.Printf("open website: %v", err)
			}
		}
	}

	runTray(trayIcon, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("[Start] server shutdown error: %v", err)
		}
	})
}
