package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Zone struct {
	ID      string
	Name    string
	Records []Record
}

func (z *Zone) Save(db *gorm.DB) error {
	result := db.Create(&z)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

type Record struct {
	Name   string `gorm:"primaryKey"`
	Type   string
	Target string
	ZoneID string
}

func OpenOrCreate(path string) (*gorm.DB, error) {
	if _, err := os.Stat(path); err == nil {
		return openDB(path)
	} else if os.IsNotExist(err) {
		return createDB(path)
	} else {
		return nil, err
	}
}

func createDB(path string) (*gorm.DB, error) {
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	db, err := openDB(path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func openDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Zone{}, &Record{})
	if err != nil {
		return nil, err
	}

	return db, err
}

func FindAllZones(db *gorm.DB) ([]Zone, error) {
	var zones []Zone

	result := db.Preload(clause.Associations).Find(&zones)
	if result.Error != nil {
		return nil, result.Error
	}

	return zones, nil
}
