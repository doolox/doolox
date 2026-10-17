package cms

import "log"

const (
	// websiteKey is the key in the key-value store holding the current website name.
	websiteKey = "website"
	// defaultWebsite is used when no website name has been stored yet.
	defaultWebsite = "doolox"
)

// Website returns the name of the current website from the key-value store.
func Website() string {
	var kv KeyValue
	if db.Where("key = ?", websiteKey).First(&kv).RecordNotFound() {
		return defaultWebsite
	}
	if kv.Value == "" {
		return defaultWebsite
	}
	return kv.Value
}

// Pages returns all pages ordered by URL.
func Pages() []Page {
	pages := []Page{}
	if err := db.Order("url asc").Find(&pages).Error; err != nil {
		log.Printf("[Pages] error: %v", err)
		return []Page{}
	}
	return pages
}

// DeletePage removes the page with the given URL.
func DeletePage(url string) error {
	return db.Where("url = ?", url).Delete(&Page{}).Error
}
