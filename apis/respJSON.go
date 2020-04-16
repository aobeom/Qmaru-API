package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DataHandler Response 数据结构
func DataHandler(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
