package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetItems(c *gin.Context)
	GetItem(c *gin.Context)
	CreateItem(c *gin.Context)
	UpdateItem(c *gin.Context)
	DeleteItem(c *gin.Context)
}

type CRUDGroup struct {
	Handler   Handler
	Namespace string
}

func StartServer(port int, segments ...CRUDGroup) (err error) {

	router := gin.Default()

	for _, v := range segments {
		group := router.Group(fmt.Sprintf("/%s", v.Namespace))
		group.GET("", v.Handler.GetItems)
		group.GET("/:id", v.Handler.GetItem)
		group.POST("", v.Handler.CreateItem)
		group.PUT("/:id", v.Handler.UpdateItem)
		group.DELETE("/:id", v.Handler.DeleteItem)
	}

	return router.Run(fmt.Sprintf(":%d", port))
}
