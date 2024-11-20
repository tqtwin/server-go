package main

import (
	"employee-management/config"
	"employee-management/middleware"
	"employee-management/models"
	"employee-management/routes"
	"fmt"

	_ "employee-management/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Employee Management API
// @version 1.0
// @description API for managing employees, departments, positions, and salaries.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 127.0.0.1:8080
// @BasePath /
func main() {
	// Kết nối đến cơ sở dữ liệu
	config.Connect()

	// Tự động migrate bảng Employee và các bảng khác
	err := config.GetDB().AutoMigrate(
		&models.Department{},
		&models.Position{},
		&models.Employee{},
		&models.Salary{},
		&models.WorkAssignment{},
		&models.EmployeeDepartment{},
		&models.EmployeePosition{},
	)
	if err != nil {
		fmt.Println("Error during migration:", err)
		return
	}
	fmt.Println("Tables migrated successfully.")

	// Khởi tạo router
	router := routes.SetupRouter()

	// Áp dụng middleware CORS
	router.Use(middleware.CORS())

	// Route Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("http://127.0.0.1:8080/swagger/doc.json"), // URL API Swagger
	))

	// Chạy server
	router.Run(":8080")
}
