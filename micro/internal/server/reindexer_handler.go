package server

import (
	"log"
	"net/http"
	"relorer/test_go_microservice/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DB interface {
	GetDocuments(limit, offset int) ([]model.Document, error)
	CreateDocument(document *model.Document) (*model.Document, error)
	GetDocument(id string) (*model.Document, error)
	UpdateDocument(document *model.Document) error
	DeleteDocument(id string) error
}

type ReindexerHandler struct {
	db DB
}

func NewHandler(db DB) *ReindexerHandler {
	return &ReindexerHandler{db: db}
}

func (h *ReindexerHandler) GetDocuments(c *gin.Context) {
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

	results, err := h.db.GetDocuments(limitInt, offsetInt)
	if err != nil {
		log.Printf("Error fetching documents: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *ReindexerHandler) CreateDocument(c *gin.Context) {
	document := &model.Document{}
	err := c.ShouldBindJSON(&document)
	if err != nil {
		log.Printf("Error binding document from JSON: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	document, err = h.db.CreateDocument(document)
	if err != nil {
		log.Printf("Error inserting document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (h *ReindexerHandler) GetDocument(c *gin.Context) {
	id := c.Param("id")

	result, err := h.db.GetDocument(id)
	if err != nil {
		log.Printf("Error fetching document: %s", err.Error())
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReindexerHandler) UpdateDocument(c *gin.Context) {
	document := model.Document{}
	err := c.ShouldBindJSON(&document)
	if err != nil {
		log.Printf("Error binding document from JSON: %s", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.db.UpdateDocument(&document)
	if err != nil {
		log.Printf("Error updating document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *ReindexerHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")

	err := h.db.DeleteDocument(id)
	if err != nil {
		log.Printf("Error deleting document: %s", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
