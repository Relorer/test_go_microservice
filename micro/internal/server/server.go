package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetDocuments(c *gin.Context)
	CreateDocument(c *gin.Context)
	GetDocument(c *gin.Context)
	UpdateDocument(c *gin.Context)
	DeleteDocument(c *gin.Context)
}

func StartServer(handler Handler, port int) (err error) {

	router := gin.Default()

	router.GET("/documents", handler.GetDocuments)
	router.POST("/documents", handler.CreateDocument)
	router.GET("/documents/:id", handler.GetDocument)
	router.PUT("/documents/:id", handler.UpdateDocument)
	router.DELETE("/documents/:id", handler.DeleteDocument)

	return router.Run(fmt.Sprintf(":%d", port))
}
