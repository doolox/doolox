package cms

import (
	"log"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", filepath.Join(dataRoot(), "doolox.db"))

	if err != nil {
		log.Printf("[initDb] error: %s", err)
	}

	db.DB()
	db.DB().Ping()
	db.LogMode(conf.Debug)
	db.AutoMigrate(&Page{}, &KeyValue{})

	// Seed the website name if it isn't configured yet.
	if db.Where("key = ?", websiteKey).First(&KeyValue{}).RecordNotFound() {
		if err := db.Create(&KeyValue{Key: websiteKey, Type: "string", Value: defaultWebsite}).Error; err != nil {
			log.Printf("[initDb] could not seed %q: %v", websiteKey, err)
		}
	}

	return db
}
