package entities

import "time"

type Category struct {
	ID        int
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
