package record

import (
	"gorm.io/gorm"
)

type Record struct {
	DB *gorm.DB
}

func NewRecord(db *gorm.DB) *Record {
	return &Record{
		DB: db,
	}
}
