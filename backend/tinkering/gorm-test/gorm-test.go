package main

import (
	"database/sql"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormTest struct {
	ID        uint
	Name      string
	Birthyear sql.NullInt16
	Invites   pq.StringArray `gorm:"type:varchar[]"`
}

func main() {
	dsn := "host=localhost user=admin password=admin dbname=swim-results port=5432 sslmode=disable TimeZone=Europe/Vienna"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&GormTest{})

	gorms := []GormTest{
		{
			ID:        14,
			Name:      "Lukas Hirsch",
			Birthyear: sql.NullInt16{Int16: 14, Valid: true},
			Invites:   []string{"Test1", "Test2"},
		},
		{
			ID:        15,
			Name:      "Lukas Hirsch",
			Birthyear: sql.NullInt16{Int16: 15, Valid: true},
		},
		{
			ID:        140,
			Name:      "Lukas Hirsch",
			Birthyear: sql.NullInt16{Int16: 140, Valid: true},
		},
	}

	// Create
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(gorms)

	// Read
	var gormTestObject GormTest
	// db.First(&gormTestObject, 1)                                       // find product with integer primary key
	db.First(&gormTestObject, GormTest{ID: 3}) // find product with code D42

	// Update - update product's price to 200
	db.Model(&gormTestObject).Update("birthyear", 2005)
	// Update - update multiple fields
	// db.Model(&gormTestObject).Updates(GormTest{Name: "Lukas"}) // non-zero fields
	// db.Model(&gormTestObject).Updates(map[string]interface{}{"Birthyear": 200, "Name": "F42"})

	// Delete - delete product
	// db.Delete(&gormTestObject, 1)
	db.Delete(&GormTest{}, 1)
}
