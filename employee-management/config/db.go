package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" // Cập nhật sử dụng PostgreSQL
	"gorm.io/gorm"
)

var DB *gorm.DB

// Kết nối đến PostgreSQL
func Connect() {
	// Load môi trường từ .env (nếu cần)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Lấy thông tin môi trường hiện tại (dev, qc, prod)
	env := os.Getenv("ENV")
	var user, password, host, port, dbName string

	// Dựa trên môi trường, thiết lập thông tin đăng nhập cơ sở dữ liệu
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

	// Tạo DSN (Data Source Name) cho PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Ho_Chi_Minh",
		host,
		user,
		password,
		dbName,
		port,
	)

	// Kết nối tới cơ sở dữ liệu PostgreSQL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("DSN: %s\n", dsn)
		panic("failed to connect to database")
	}

	log.Println("Database connected successfully!")
}

// Lấy kết nối DB
func GetDB() *gorm.DB {
	return DB
}
