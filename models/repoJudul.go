package models

import "github.com/lib/pq"

type Title struct {
	Titles pq.StringArray `gorm:"type:text[]" json:"title"`
}
