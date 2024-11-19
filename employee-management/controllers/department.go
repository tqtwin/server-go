package controllers

import (
	"employee-management/config"
	"employee-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDepartments godoc
// @Summary Get all departments
// @Description Retrieve a list of all departments
// @Tags Department
// @Accept json
// @Produce json
// @Success 200 {array} models.Department
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments [get]
func GetDepartments(c *gin.Context) {
	var departments []models.Department
	if err := config.GetDB().Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch departments"})
		return
	}
	c.JSON(http.StatusOK, departments)
}

// GetEmployeesByDepartment godoc
// @Summary Get employees by department
// @Description Retrieve all employees in a specific department
// @Tags Department
// @Accept json
// @Produce json
// @Param department_id path int true "Department ID"
// @Success 200 {array} models.Employee
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{department_id}/employees [get]
func GetEmployeesByDepartment(c *gin.Context) {
	departmentID := c.Param("department_id")

	var employees []models.Employee

	// Assuming employee_departments is a junction table
	if err := config.GetDB().Table("employee_departments").
		Where("employee_departments.department_id = ?", departmentID).
		Joins("JOIN employees ON employee_departments.employee_id = employees.id").
		Select("employees.id, employees.name").
		Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch employees for this department"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// CreateDepartment godoc
// @Summary Create a new department
// @Description Add a new department record to the database
// @Tags Department
// @Accept json
// @Produce json
// @Param department body models.Department true "Department data"
// @Success 201 {object} models.Department
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments [post]
func CreateDepartment(c *gin.Context) {
	var department models.Department

	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Check if the department already exists
	var existingDepartment models.Department
	if err := config.GetDB().Where("name = ?", department.Name).First(&existingDepartment).Error; err == nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Department with this name already exists"})
		return
	}

	// Create the new department
	if err := config.GetDB().Create(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, department)
}

// UpdateDepartment godoc
// @Summary Update a department
// @Description Update an existing department's information
// @Tags Department
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param department body models.Department true "Updated department data"
// @Success 200 {object} models.Department
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{id} [put]
func UpdateDepartment(c *gin.Context) {
	var department models.Department
	id := c.Param("id")

	// Find the department by ID
	if err := config.GetDB().Where("id = ?", id).First(&department).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Department not found"})
		return
	}

	if err := c.ShouldBindJSON(&department); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	if err := config.GetDB().Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// DeleteDepartment godoc
// @Summary Delete a department
// @Description Remove a department record from the database
// @Tags Department
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/departments/{id} [delete]
func DeleteDepartment(c *gin.Context) {
	var department models.Department
	id := c.Param("id")

	// Find the department by ID
	if err := config.GetDB().Where("id = ?", id).First(&department).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Department not found"})
		return
	}

	// Check if any employees are assigned to the department
	var employees []models.Employee
	if err := config.GetDB().Where("department_id = ?", id).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to check employees in department"})
		return
	}
	if len(employees) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Cannot delete department because it has employees"})
		return
	}

	// Delete the department
	if err := config.GetDB().Delete(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete department"})
		return
	}

	c.JSON(http.StatusOK, ResponseMessage{Message: "Department deleted successfully"})
}
