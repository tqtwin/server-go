package middleware

import "github.com/gin-gonic/gin"

// CORS Middleware: Cho phép xử lý các yêu cầu từ nguồn khác
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy giá trị của "Origin" từ header yêu cầu
		origin := c.Request.Header.Get("Origin")

		// Thiết lập các header CORS
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin) // Chỉ định nguồn gốc
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // Cho phép sử dụng thông tin xác thực (cookie, headers, etc.)

		// Nếu là yêu cầu OPTIONS (preflight request), phản hồi ngay lập tức
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Tiếp tục đến handler tiếp theo
		c.Next()
	}
}
