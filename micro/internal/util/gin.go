package util

import (
	"log"

	"github.com/gin-gonic/gin"
)

func GinHandleError(c *gin.Context, err error, status int) bool {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		c.AbortWithStatus(status)
		return true
	}
	return false
}
