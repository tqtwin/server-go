package controllers

import (
	"employee-management/config"
	"employee-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSalaries godoc
// @Summary Get list of salaries with filters
// @Description Get salaries with optional filters for month, quarter, year, and payment status
// @Tags Salary
// @Accept json
// @Produce json
// @Param month query int false "Month"
// @Param quarter query int false "Quarter"
// @Param year query int false "Year"
// @Param status query string false "Status (Chưa thanh toán/Đã thanh toán)"
// @Success 200 {array} models.Salary
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/salaries [get]
func GetSalaries(c *gin.Context) {
	// Get query parameters from the URL
	month := c.DefaultQuery("month", "")
	quarter := c.DefaultQuery("quarter", "")
	year := c.DefaultQuery("year", "")
	status := c.DefaultQuery("status", "")

	var salaries []models.Salary
	query := config.GetDB().Model(&models.Salary{})

	// Handle 'month' filter
	if month != "" {
		monthInt, err := strconv.Atoi(month)
		if err == nil {
			query = query.Where("MONTH(created_at) = ?", monthInt)
		}
	}

	// Handle 'quarter' filter
	if quarter != "" {
		quarterInt, err := strconv.Atoi(quarter)
		if err == nil {
			query = query.Where("QUARTER(created_at) = ?", quarterInt)
		}
	}

	// Handle 'year' filter
	if year != "" {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			query = query.Where("YEAR(created_at) = ?", yearInt)
		}
	}

	// Handle 'status' filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Execute the query to fetch salaries
	if err := query.Find(&salaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch salaries"})
		return
	}

	// Return the filtered results
	c.JSON(http.StatusOK, salaries)
}

// GetSalaryByID godoc
// @Summary Get salary by ID
// @Description Retrieve a salary by its ID, including employee information
// @Tags Salary
// @Accept json
// @Produce json
// @Param id path int true "Salary ID"
// @Success 200 {object} models.Salary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/salaries/{id} [get]
func GetSalaryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid salary ID"})
		return
	}

	var salary models.Salary
	if err := config.GetDB().Preload("Employee").First(&salary, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Salary not found"})
		return
	}
	c.JSON(http.StatusOK, salary)
}

// CreateSalary godoc
// @Summary Create a new salary
// @Description Add a new salary record for an employee
// @Tags Salary
// @Accept json
// @Produce json
// @Param salary body models.Salary true "Salary data"
// @Success 201 {object} models.Salary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/salaries [post]
// CreateSalary godoc
func CreateSalary(c *gin.Context) {
	var salary models.Salary
	if err := c.ShouldBindJSON(&salary); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Kiểm tra xem EmployeeID có tồn tại trong hệ thống không
	var employee models.Employee
	if err := config.GetDB().First(&employee, salary.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Employee not found"})
		return
	}

	// Thiết lập trạng thái mặc định là "Chưa thanh toán"
	salary.Status = "Chưa thanh toán"

	// Lưu bảng lương vào cơ sở dữ liệu
	if err := config.GetDB().Create(&salary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create salary"})
		return
	}

	c.JSON(http.StatusCreated, salary)
}

// UpdateSalary godoc
// @Summary Update an existing salary
// @Description Update a salary record by its ID
// @Tags Salary
// @Accept json
// @Produce json
// @Param id path int true "Salary ID"
// @Param salary body models.Salary true "Updated salary data"
// @Success 200 {object} models.Salary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/salaries/{id} [put]
// UpdateSalary godoc
func UpdateSalary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid salary ID"})
		return
	}

	var salary models.Salary
	if err := config.GetDB().First(&salary, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Salary not found"})
		return
	}

	if err := c.ShouldBindJSON(&salary); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Cập nhật bảng lương
	if err := config.GetDB().Save(&salary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update salary"})
		return
	}
	c.JSON(http.StatusOK, salary)
}

// DeleteSalary godoc
// @Summary Delete a salary
// @Description Remove a salary record by its ID
// @Tags Salary
// @Accept json
// @Produce json
// @Param id path int true "Salary ID"
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/salaries/{id} [delete]
func DeleteSalary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid salary ID"})
		return
	}

	var salary models.Salary
	if err := config.GetDB().First(&salary, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Salary not found"})
		return
	}

	// Xóa bảng lương khỏi cơ sở dữ liệu
	if err := config.GetDB().Delete(&salary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete salary"})
		return
	}
	c.JSON(http.StatusOK, ResponseMessage{Message: "Salary deleted successfully"})
}

// PaySalary godoc
// @Summary Mark a salary as paid
// @Description Mark a salary record as paid by changing the status to "Đã thanh toán"
// @Tags Salary
// @Accept json
// @Produce json
// @Param id path int true "Salary ID"
// @Success 200 {object} models.Salary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/salaries/{id}/pay [put]
func PaySalary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid salary ID"})
		return
	}

	var salary models.Salary
	if err := config.GetDB().First(&salary, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Salary not found"})
		return
	}

	// Cập nhật trạng thái thành "Đã thanh toán"
	salary.Status = "Đã thanh toán"
	if err := config.GetDB().Save(&salary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update salary"})
		return
	}

	c.JSON(http.StatusOK, salary)
}

// GetSalaryStatistics godoc
// @Summary Get salary statistics
// @Description Get total salaries by filters: month, year, and payment status
// @Tags Salary
// @Accept json
// @Produce json
// @Param month query int false "Month"
// @Param year query int false "Year"
// @Param status query string false "Status (Chưa thanh toán/Đã thanh toán)"
// @Success 200 {object} map[string]float64 "Statistics"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/salaries/stats [get]
func GetSalaryStatistics(c *gin.Context) {
	// Get query parameters for filters
	month := c.DefaultQuery("month", "")
	year := c.DefaultQuery("year", "")
	status := c.DefaultQuery("status", "")

	var totalSalaries float64
	query := config.GetDB().Model(&models.Salary{})

	// Apply filters based on query parameters
	if month != "" {
		monthInt, err := strconv.Atoi(month)
		if err == nil {
			query = query.Where("MONTH(created_at) = ?", monthInt)
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid month format"})
			return
		}
	}
	if year != "" {
		yearInt, err := strconv.Atoi(year)
		if err == nil {
			query = query.Where("YEAR(created_at) = ?", yearInt)
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid year format"})
			return
		}
	}
	if status != "" {
		// Ensure status is either 'Chưa thanh toán' or 'Đã thanh toán'
		if status != "Chưa thanh toán" && status != "Đã thanh toán" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid status. Must be 'Chưa thanh toán' or 'Đã thanh toán'"})
			return
		}
		query = query.Where("status = ?", status)
	}

	// Calculate total salaries with the applied filters
	if err := query.Select("SUM(total_salary)").Scan(&totalSalaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to calculate salary statistics"})
		return
	}

	// Return the result as JSON response
	c.JSON(http.StatusOK, gin.H{"total": totalSalaries})
}
