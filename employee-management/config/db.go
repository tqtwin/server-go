package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the current environment (dev, qc, prod)
	env := os.Getenv("ENV")
	var user, password, host, port, dbName string

	// Based on the environment, set the database credentials
	switch env {
	case "dev":
		user = os.Getenv("DEV_DB_USER")
		password = os.Getenv("DEV_DB_PASSWORD")
		host = os.Getenv("DEV_DB_HOST")
		port = os.Getenv("DEV_DB_PORT")
		dbName = os.Getenv("DEV_DB_NAME")
	case "qc":
		user = os.Getenv("QC_DB_USER")
		password = os.Getenv("QC_DB_PASSWORD")
		host = os.Getenv("QC_DB_HOST")
		port = os.Getenv("QC_DB_PORT")
		dbName = os.Getenv("QC_DB_NAME")
	case "prod":
		user = os.Getenv("PROD_DB_USER")
		password = os.Getenv("PROD_DB_PASSWORD")
		host = os.Getenv("PROD_DB_HOST")
		port = os.Getenv("PROD_DB_PORT")
		dbName = os.Getenv("PROD_DB_NAME")
	default:
		log.Fatalf("Unknown environment: %s", env)
	}

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, password, host, port, dbName)

	// Connect to the database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf(dsn)
		panic("failed to connect database")
	}
}

func GetDB() *gorm.DB {
	return DB
}
