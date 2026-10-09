package cms

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/macaron.v1"
)

var m *macaron.Macaron

var conf *Config

var db *gorm.DB

func init() {
	conf = initConfig()

	db = initDb()

	m = initMacaron()
}

func Start() {
	m.Get("/", pageView)
	m.Get("/generate", generateView)
	m.Get("/:page", pageView)

	m.Run("0.0.0.0", 5000)
}
