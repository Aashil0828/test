package models

import "gorm.io/gorm"

type AddressJoin struct {
	gorm.Model
	PersonID  uint `json:"person_id"`
	AddressID uint `json:"address_id"`
}
