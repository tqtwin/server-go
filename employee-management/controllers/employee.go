// @BasePath /api/v1
package controllers

import (
	"employee-management/config"
	"employee-management/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ErrorResponse represents a generic error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// ResponseMessage represents a generic success message
type ResponseMessage struct {
	Message string `json:"message"`
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if the given password matches the hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// isValidEmail checks the validity of an email address
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

// RegisterEmployee godoc
// @Summary Register a new employee
// @Description Creates a new employee record with hashed password
// @Tags Employee
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee data"
// @Success 201 {object} models.Employee
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/employees/register [post]
func RegisterEmployee(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Check if email is already in use
	var existingEmployee models.Employee
	if err := config.GetDB().Where("email = ?", employee.Email).First(&existingEmployee).Error; err == nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Email already registered"})
		return
	}

	// Validate email format
	if !isValidEmail(employee.Email) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email format"})
		return
	}

	// Validate password length
	if len(employee.Password) < 6 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Password must be at least 6 characters"})
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(employee.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}
	employee.Password = hashedPassword

	// Save employee record
	if err := config.GetDB().Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// LoginEmployee godoc
// @Summary Login an employee
// @Description Authenticate an employee using email and password
// @Tags Employee
// @Accept json
// @Produce json
// @Param loginData body map[string]string true "Email and Password"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/employees/login [post]
type LoginResponse struct {
	Message string `json:"message"`
	IsAdmin bool   `json:"isAdmin"`
}

func LoginEmployee(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email and password are required"})
		return
	}

	var employee models.Employee
	if err := config.GetDB().Where("email = ?", loginData.Email).Preload("EmployeeDepartments").Preload("EmployeePositions").First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Employee not found"})
		return
	}

	// Validate password
	if !CheckPasswordHash(loginData.Password, employee.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Determine if the user is an admin
	isAdmin := employee.Role == "Admin"

	// Extract department IDs and position IDs
	var departmentIDs []uint
	var positionIDs []uint

	for _, ed := range employee.EmployeeDepartments {
		departmentIDs = append(departmentIDs, ed.DepartmentID)
	}

	for _, ep := range employee.EmployeePositions {
		positionIDs = append(positionIDs, ep.PositionID)
	}

	// Respond with full employee data
	c.JSON(http.StatusOK, gin.H{
		"message":        "Login successful",
		"isAdmin":        isAdmin,
		"id":             employee.ID,
		"name":           employee.Name,
		"email":          employee.Email,
		"phone":          employee.Phone,
		"address":        employee.Address,
		"gender":         employee.Gender,
		"date_of_birth":  employee.DateOfBirth,
		"status":         employee.Status,
		"role":           employee.Role,
		"department_ids": departmentIDs,
		"position_ids":   positionIDs,
		"created_at":     employee.CreatedAt,
		"updated_at":     employee.UpdatedAt,
	})
}

// GetEmployees godoc
// @Summary Get all employees
// @Description Retrieve a list of all employees
// @Tags Employee
// @Accept json
// @Produce json
// @Success 200 {array} models.Employee
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/employees [get]
func GetEmployees(c *gin.Context) {
	var employees []models.Employee

	// Sử dụng Preload để tải thông tin các department và position của mỗi nhân viên
	if err := config.GetDB().Preload("EmployeeDepartments").Preload("EmployeePositions").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch employees"})
		return
	}

	// Lấy danh sách department_ids và position_ids từ các quan hệ đã preload
	for i, employee := range employees {
		// Lấy danh sách department_ids từ EmployeeDepartments
		for _, ed := range employee.EmployeeDepartments {
			employees[i].DepartmentIDs = append(employees[i].DepartmentIDs, ed.DepartmentID)
		}

		// Lấy danh sách position_ids từ EmployeePositions
		for _, ep := range employee.EmployeePositions {
			employees[i].PositionIDs = append(employees[i].PositionIDs, ep.PositionID)
		}
	}

	// Trả về thông tin nhân viên
	c.JSON(http.StatusOK, employees)
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Add a new employee record to the database
// @Tags Employee
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee data"
// @Success 201 {object} models.Employee
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/employees [post]
func CreateEmployee(c *gin.Context) {
	var employee models.Employee

	// Bind dữ liệu JSON vào struct Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid input data: %v", err)})
		return
	}

	// Insert the main employee record first to get employee.ID
	if err := config.GetDB().Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create employee"})
		return
	}

	// Validate and handle department IDs
	if len(employee.DepartmentIDs) > 0 {
		for _, deptID := range employee.DepartmentIDs {
			var department models.Department
			if err := config.GetDB().Where("id = ?", deptID).First(&department).Error; err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid department ID"})
				return
			}
			// Insert into employee_departments table
			if err := config.GetDB().Create(&models.EmployeeDepartment{
				EmployeeID:   employee.ID,
				DepartmentID: deptID,
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add employee to department"})
				return
			}
		}
	}

	// Validate and handle position IDs
	if len(employee.PositionIDs) > 0 {
		for _, posID := range employee.PositionIDs {
			var position models.Position
			if err := config.GetDB().Where("id = ?", posID).First(&position).Error; err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid position ID"})
				return
			}
			// Insert into employee_positions table
			if err := config.GetDB().Create(&models.EmployeePosition{
				EmployeeID: employee.ID,
				PositionID: posID,
			}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add employee to position"})
				return
			}
		}
	}

	// Trả về thông tin nhân viên đã tạo
	c.JSON(http.StatusCreated, employee)
}

// UpdateEmployee godoc
// @Summary Update an employee
// @Description Update details of an existing employee, including multiple departments and positions
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body models.Employee true "Updated employee data"
// @Success 200 {object} models.Employee
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/employees/{id} [put]
func UpdateEmployee(c *gin.Context) {
	var employee models.Employee
	id := c.Param("id")

	// Check if employee exists
	if err := config.GetDB().Where("id = ?", id).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Employee not found"})
		return
	}

	// Bind updated data
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Validate departments and positions as done in CreateEmployee
	// Same logic for department and position validation as in CreateEmployee function...

	// Update employee record
	if err := config.GetDB().Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// DeleteEmployee godoc
// @Summary Delete an employee
// @Description Remove an employee record from the database
// @Tags Employee
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} ResponseMessage
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/employees/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	var employee models.Employee
	id := c.Param("id")

	// Check if employee exists
	if err := config.GetDB().Where("id = ?", id).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Employee not found"})
		return
	}

	// Xóa các bản ghi liên quan trong employee_departments và employee_positions
	if err := config.GetDB().Where("employee_id = ?", id).Delete(&models.EmployeeDepartment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete employee departments"})
		return
	}

	if err := config.GetDB().Where("employee_id = ?", id).Delete(&models.EmployeePosition{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete employee positions"})
		return
	}

	// Delete the employee record
	if err := config.GetDB().Delete(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, ResponseMessage{Message: "Employee deleted successfully"})
}
