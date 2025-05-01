package entities

import "time"

type Product struct {
	ID          int
	Name        string
	Price       float64
	CategoryID  uint
	Stock       int
	Description string
	Category    Category
	UpdatedAt   time.Time
	CreatedAt   time.Time
}
