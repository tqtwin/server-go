package controllers

import (
	"employee-management/config"
	"employee-management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Lấy danh sách công tác
func GetWorkAssignments(c *gin.Context) {
	var workAssignments []models.WorkAssignment
	if err := config.GetDB().Preload("Employee").Find(&workAssignments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch work assignments"})
		return
	}
	c.JSON(http.StatusOK, workAssignments)
}

// Thêm một công tác mới
func CreateWorkAssignment(c *gin.Context) {
	var request models.WorkAssignment

	// Giải mã JSON từ request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tìm kiếm nhân viên theo EmployeeID
	var employee models.Employee
	if err := config.GetDB().First(&employee, request.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Gán EmployeeName từ nhân viên và đảm bảo Status có giá trị mặc định
	request.EmployeeName = employee.Name
	if request.Status == "" {
		request.Status = "Chưa hoàn thành"
	}

	// Lưu công việc mới vào cơ sở dữ liệu
	if err := config.GetDB().Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create work assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work assignment created successfully", "data": request})
}

// Cập nhật công tác
func UpdateWorkAssignment(c *gin.Context) {
	var workAssignment models.WorkAssignment
	id := c.Param("id")

	// Tìm công tác theo ID
	if err := config.GetDB().Where("id = ?", id).First(&workAssignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work assignment not found"})
		return
	}

	// Liên kết dữ liệu JSON mới với công tác
	if err := c.ShouldBindJSON(&workAssignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Lưu thông tin công tác đã cập nhật vào cơ sở dữ liệu
	if err := config.GetDB().Save(&workAssignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work assignment"})
		return
	}

	c.JSON(http.StatusOK, workAssignment)
}

// Xóa công tác
func DeleteWorkAssignment(c *gin.Context) {
	var workAssignment models.WorkAssignment
	id := c.Param("id")

	// Tìm công tác theo ID
	if err := config.GetDB().Where("id = ?", id).First(&workAssignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Work assignment not found"})
		return
	}

	// Xóa công tác khỏi cơ sở dữ liệu
	if err := config.GetDB().Delete(&workAssignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete work assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work assignment deleted successfully"})
}
