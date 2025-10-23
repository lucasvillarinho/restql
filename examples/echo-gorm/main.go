package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global GORM database connection
var db *gorm.DB

func main() {
	// Initialize GORM database (SQLite in-memory)
	initDB()

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes
	setupRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// initDB initializes the GORM SQLite in-memory database with example data
func initDB() {
	var err error

	// Initialize GORM with SQLite
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// Seed database with example data
	if err := seedDatabase(); err != nil {
		panic("failed to seed database: " + err.Error())
	}
}
