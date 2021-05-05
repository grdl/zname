package zname

import (
	"database/sql"
	"errors"
	"fmt"
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

func DeleteCache(path string) error {
	err := os.Remove(path)

	// Don't throw errors when trying to delete a non-existent file.
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func OpenCache(path string) (*gorm.DB, error) {
	if _, err := os.Stat(path); err == nil {
		return openDB(path)
	} else if os.IsNotExist(err) {
		return createDB(path)
	} else {
		return nil, err
	}
}

func RebuildCache(path string) error {
	client, err := NewFromConfig()
	if err != nil {
		return err
	}

	err = DeleteCache(path)
	if err != nil {
		return err
	}

	db, err := OpenCache(path)
	if err != nil {
		return err
	}

	zones, err := client.GetZones()
	if err != nil {
		return err
	}

	for _, zone := range zones {
		fmt.Printf("Scraping %s zone...\n", zone.Name)

		records, err := client.GetRecords(zone.ID)
		if err != nil {
			return err
		}

		fmt.Printf("\tFound %d records\n", len(records))
		zone.Records = records

		err = zone.Save(db)
		if err != nil {
			return err
		}
	}

	return nil
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

func FindByWord(db *gorm.DB, word string) ([]Record, error) {
	var records []Record

	result := db.Where("name LIKE @word OR target LIKE @word", sql.Named("word", fmt.Sprintf("%%%s%%", word))).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}
