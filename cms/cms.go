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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	setupDataDir()

	conf = initConfig()

	setupWebsiteDir()

	m = initMacaron()
}

// routes registers all HTTP routes on the given macaron instance.
func routes(m *macaron.Macaron) {
	m.Get("/generate", generateView)
	m.Get("/publish", publishView)
	m.Get("/doolox.js", dooloxJSView)
	m.Get("/", homeView)
	m.Get("/assets/*", assetsView)
	m.Get("/:seg", segmentView)
	m.Get("/:seg/:page", nestedView)
	m.Get("/:seg/assets/*", nestedAssetsView)
}

// Start runs the Doolox web server and shows the system tray icon. It blocks
// until the tray is quit, then shuts the server down and returns.
func Start(trayIcon TrayIcon, openWebsite bool) {
	routes(m)

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
