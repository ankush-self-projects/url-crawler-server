package model

import "time"

type URL struct {
	ID            uint `gorm:"primaryKey"`
	URL           string
	HTMLVersion   string
	PageTitle     string
	Headings      string
	InternalLinks int
	ExternalLinks int
	BrokenLinks   int
	HasLoginForm  bool
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
