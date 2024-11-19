package routes

import (
	"employee-management/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Group API v1
	apiV1 := router.Group("/api/v1")
	{
		// Routes cho Employee
		employeeRoutes := apiV1.Group("/employees")
		{
			employeeRoutes.POST("/login", controllers.LoginEmployee)
			employeeRoutes.POST("/register", controllers.RegisterEmployee)
			employeeRoutes.GET("/", controllers.GetEmployees)
			employeeRoutes.POST("/", controllers.CreateEmployee)
			employeeRoutes.PUT("/:id", controllers.UpdateEmployee)
			employeeRoutes.DELETE("/:id", controllers.DeleteEmployee)
		}

		// Routes cho Department
		departmentRoutes := apiV1.Group("/departments")
		{
			departmentRoutes.GET("/", controllers.GetDepartments)
			departmentRoutes.GET("/:department_id/employees", controllers.GetEmployeesByDepartment)
			departmentRoutes.POST("/", controllers.CreateDepartment)
			departmentRoutes.PUT("/:id", controllers.UpdateDepartment)
			departmentRoutes.DELETE("/:id", controllers.DeleteDepartment)
		}

		// Routes cho Position
		positionRoutes := apiV1.Group("/positions")
		{
			positionRoutes.GET("/", controllers.GetPositions)
			// positionRoutes.GET("/:id", controllers.GetPositionByID)
			positionRoutes.GET("/:position_id/employees", controllers.GetEmployeesByPosition)
			positionRoutes.POST("/", controllers.CreatePosition)
			positionRoutes.PUT("/:id", controllers.UpdatePosition)
			positionRoutes.DELETE("/:id", controllers.DeletePosition)
		}
		salaries := apiV1.Group("/salaries")
		{
			salaries.GET("/", controllers.GetSalaries)      // Lấy tất cả bảng lương
			salaries.GET("/:id", controllers.GetSalaryByID) // Lấy bảng lương theo ID
			salaries.POST("/", controllers.CreateSalary)    // Tạo bảng lương mới
			salaries.PUT("/:id", controllers.UpdateSalary)  // Cập nhật bảng lương theo ID
			salaries.DELETE("/:id", controllers.DeleteSalary)
			salaries.PUT("/:id/pay", controllers.PaySalary) // Xóa bảng lương theo ID
			salaries.GET("/stats", controllers.GetSalaryStatistics)

		}
		workassignments := apiV1.Group("/workassignments")
		{
			workassignments.GET("/", controllers.GetWorkAssignments)
			workassignments.POST("/", controllers.CreateWorkAssignment)
			workassignments.PUT("/:id", controllers.UpdateWorkAssignment)
			workassignments.DELETE("/:id", controllers.DeleteWorkAssignment)
		}
	}

	return router
}
