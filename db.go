package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// DB Internal database structure
type DB struct {
	db gorm.DB
}

// SetupDb Connect with postgres
func SetupDb() (DB, error) {
	db, err := gorm.Open("postgres", "sslmode=disable")

	return DB{db}, err
}
