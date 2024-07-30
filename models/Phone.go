package models

import "gorm.io/gorm"

type Phone struct {
	gorm.Model
	Number   string `json:"number"`
	PersonID uint   `json:"person_id"`
}
