package cms

import "log"

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
