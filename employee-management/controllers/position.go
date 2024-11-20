package controllers

import (
	"employee-management/config"
	"employee-management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// positionExists checks if a position already exists by its title.
func positionExists(title string) bool {
	var position models.Position
	if result := config.GetDB().Where("title = ?", title).First(&position); result.Error == nil {
		return true
	}
	return false
}

// GetPositions godoc
// @Summary Get all positions
// @Description Retrieve a list of all positions
// @Tags Position
// @Accept json
// @Produce json
// @Success 200 {array} models.Position
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions [get]
func GetPositions(c *gin.Context) {
	var positions []models.Position
	if err := config.GetDB().Find(&positions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch positions"})
		return
	}
	c.JSON(http.StatusOK, positions)
}

// GetPositionByID godoc
// @Summary Get position by ID
// @Description Retrieve a position by its ID
// @Tags Position
// @Accept json
// @Produce json
// @Param id path int true "Position ID"
// @Success 200 {object} models.Position
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions/{id} [get]
func GetPositionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid position ID"})
		return
	}

	var position models.Position
	if err := config.GetDB().First(&position, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Position not found"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// GetEmployeesByPosition godoc
// @Summary Get employees by position
// @Description Retrieve employees by their position
// @Tags Position
// @Accept json
// @Produce json
// @Param position_id path int true "Position ID"
// @Success 200 {array} models.Employee
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions/{position_id}/employees [get]
func GetEmployeesByPosition(c *gin.Context) {
	positionID := c.Param("position_id")

	var employees []models.Employee

	// Query employees with the specified position ID
	if err := config.GetDB().Where("position_id = ?", positionID).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch employees for this position"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// CreatePosition godoc
// @Summary Create a new position
// @Description Add a new position to the system
// @Tags Position
// @Accept json
// @Produce json
// @Param position body models.Position true "Position data"
// @Success 201 {object} models.Position
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions [post]
func CreatePosition(c *gin.Context) {
	var position models.Position

	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input data"})
		return
	}

	// Check if the position already exists
	if positionExists(position.Title) {
		c.JSON(http.StatusConflict, models.ErrorResponse{Error: "Position already exists"})
		return
	}

	// Create the new position
	if err := config.GetDB().Create(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create position"})
		return
	}

	c.JSON(http.StatusCreated, position)
}

// UpdatePosition godoc
// @Summary Update an existing position
// @Description Modify a position's data
// @Tags Position
// @Accept json
// @Produce json
// @Param id path int true "Position ID"
// @Param position body models.Position true "Updated position data"
// @Success 200 {object} models.Position
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions/{id} [put]
func UpdatePosition(c *gin.Context) {
	var position models.Position
	id := c.Param("id")

	// Check if position exists
	if err := config.GetDB().Where("id = ?", id).First(&position).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Position not found"})
		return
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&position); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid input data"})
		return
	}

	if err := config.GetDB().Save(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update position"})
		return
	}

	c.JSON(http.StatusOK, position)
}

// DeletePosition godoc
// @Summary Delete a position
// @Description Remove a position from the system
// @Tags Position
// @Accept json
// @Produce json
// @Param id path int true "Position ID"
// @Success 200 {object} models.ResponseMessage
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/positions/{id} [delete]
func DeletePosition(c *gin.Context) {
	id := c.Param("id")

	// Find the position by ID
	var position models.Position
	if err := config.GetDB().Where("id = ?", id).First(&position).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Position not found"})
		return
	}

	// Check if any employees are assigned to the position
	var employees []models.Employee
	if err := config.GetDB().Where("position_id = ?", id).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to check employees in position"})
		return
	}
	if len(employees) > 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Cannot delete position because it has employees"})
		return
	}

	// Delete the position
	if err := config.GetDB().Delete(&position).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete position"})
		return
	}

	c.JSON(http.StatusOK, models.ResponseMessage{Message: "Position deleted successfully"})
}
