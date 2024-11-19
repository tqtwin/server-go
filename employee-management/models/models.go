package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// ResponseMessage đại diện cho phản hồi chung
type ResponseMessage struct {
	Message string `json:"message"`
}

// ErrorResponse đại diện cho lỗi
type ErrorResponse struct {
	Error string `json:"error"`
}

// CustomTime là kiểu thời gian tùy chỉnh để chỉ nhận dạng theo định dạng "YYYY-MM-DD"
type CustomTime struct {
	time.Time
}

// Implement phương thức UnmarshalJSON cho CustomTime để parse chuỗi theo định dạng "YYYY-MM-DD"
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(`"2006-01-02"`, s) // Định dạng ngày "YYYY-MM-DD"
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}
	ct.Time = t
	return nil
}

// Implement Valuer interface để lưu CustomTime vào cơ sở dữ liệu
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Implement Scanner interface để lấy CustomTime từ cơ sở dữ liệu
func (ct *CustomTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into CustomTime", value)
	}
	ct.Time = t
	return nil
}

// Employee đại diện cho một nhân viên
type Employee struct {
	ID                  uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                string               `json:"name" gorm:"not null"`
	Email               string               `json:"email" gorm:"unique;not null"`
	Password            string               `json:"password" gorm:"null"`
	Cmnd                string               `json:"cmnd" gorm:"unique;not null"`
	DateOfBirth         CustomTime           `json:"date_of_birth"`
	Phone               string               `json:"phone" gorm:"unique;not null"`
	Address             string               `json:"address"`
	Role                string               `json:"role" gorm:"not null"`
	Status              string               `json:"status" gorm:"not null"`
	Gender              string               `json:"gender"`
	CreatedAt           time.Time            `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time            `json:"updated_at" gorm:"autoUpdateTime"`
	DepartmentIDs       []uint               `json:"department_ids" gorm:"-"` // Không lưu vào database
	PositionIDs         []uint               `json:"position_ids" gorm:"-"`   // Không lưu vào database
	EmployeeDepartments []EmployeeDepartment `json:"employee_departments" gorm:"foreignKey:EmployeeID"`
	EmployeePositions   []EmployeePosition   `json:"employee_positions" gorm:"foreignKey:EmployeeID"`
}

// Department đại diện cho phòng ban trong tổ chức
type Department struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Employees   []Employee `gorm:"many2many:employee_departments" json:"employees"`
}

// Position đại diện cho một chức vụ trong tổ chức
type Position struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Employees   []Employee `gorm:"many2many:employee_positions" json:"employees"`
}

// Salary đại diện cho bảng lương của nhân viên
type Salary struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	EmployeeID   uint      `json:"employee_id" gorm:"not null"`
	EmployeeName string    `json:"employee_name" gorm:"not null"`
	BasicSalary  int       `json:"basic_salary" gorm:"not null"`
	Coefficient  int       `json:"coefficient" gorm:"not null"`
	Bonus        int       `json:"bonus" gorm:"not null"`
	Fine         int       `json:"fine" gorm:"not null"`
	TotalSalary  int       `json:"total_salary" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	WorkingDays  uint      `json:"working_days" gorm:"default:0"`
	Employee     Employee  `json:"employee" gorm:"foreignKey:EmployeeID"`
	Status       string    `json:"status" gorm:"not null;default:'unpaid'"`
}

// EmployeeDepartment đại diện cho quan hệ giữa nhân viên và phòng ban
type EmployeeDepartment struct {
	EmployeeID   uint `json:"employee_id"`
	DepartmentID uint `json:"department_id"`
}

// EmployeePosition đại diện cho quan hệ giữa nhân viên và chức vụ
type EmployeePosition struct {
	EmployeeID uint `json:"employee_id"`
	PositionID uint `json:"position_id"`
}

// WorkAssignment đại diện cho bảng công tác của nhân viên
type WorkAssignment struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	EmployeeID   uint       `json:"employee_id" gorm:"not null"`
	EmployeeName string     `json:"employee_name" gorm:"not null"`
	Assignment   string     `json:"assignment" gorm:"not null"`
	StartDate    CustomTime `json:"start_date" gorm:"not null"`
	EndDate      CustomTime `json:"end_date"`
	Status       string     `json:"status" gorm:"not null;default:'Chưa hoàn thành'"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	Employee     Employee   `json:"employee" gorm:"foreignKey:EmployeeID"`
}
