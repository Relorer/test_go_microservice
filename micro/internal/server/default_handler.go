package server

import (
	"log"
	"net/http"
	"relorer/test_go_microservice/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetDocuments(limit, offset int) ([]*model.Document, error)
	CreateDocument(document *model.Document) (*model.Document, error)
	GetDocument(id int64) (*model.Document, error)
	UpdateDocument(document *model.Document) error
	DeleteDocument(id int64) error
}

type DefaultHandler struct {
	repo Repository
}

func NewHandler(repo Repository) *DefaultHandler {
	return &DefaultHandler{repo: repo}
}

func (h *DefaultHandler) GetDocuments(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("Error converting limit to integer: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("Error converting offset to integer: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	results, err := h.repo.GetDocuments(limitInt, offsetInt)
	if err != nil {
		log.Printf("Error fetching documents: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DefaultHandler) CreateDocument(c *gin.Context) {
	document := &model.Document{}
	err := c.ShouldBindJSON(&document)
	if err != nil {
		log.Printf("Error binding document from JSON: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Println(document)
	document, err = h.repo.CreateDocument(document)
	if err != nil {
		log.Printf("Error inserting document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println(document)
	c.JSON(http.StatusCreated, document)
}

func (h *DefaultHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error fetching document: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := h.repo.GetDocument(idInt)
	if err != nil {
		log.Printf("Error fetching document: %s", err.Error())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DefaultHandler) UpdateDocument(c *gin.Context) {
	document := model.Document{}
	err := c.ShouldBindJSON(&document)
	if err != nil {
		log.Printf("Error binding document from JSON: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id := c.Param("id")

	document.ID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error updating document: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateDocument(&document)
	if err != nil {
		log.Printf("Error updating document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *DefaultHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Printf("Error deleting document: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteDocument(idInt)
	if err != nil {
		log.Printf("Error deleting document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
