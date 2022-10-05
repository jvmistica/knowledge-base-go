package record

import (
	"gorm.io/gorm"
)

// Record is the record structure
type Record struct {
	DB *gorm.DB
}

// NewRecord returns a record
func NewRecord(db *gorm.DB) *Record {
	return &Record{
		DB: db,
	}
}
