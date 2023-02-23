package handler

import (
	"net/http"
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DocumentRepository interface {
	GetDocuments(limit, offset int, join bool) ([]*model.Document, error)
	CreateDocument(document *model.Document) (*model.Document, error)
	GetDocument(id int64) (*model.Document, error)
	UpdateDocument(document *model.Document) error
	DeleteDocument(id int64) error
}

type DocumentHandler struct {
	repo DocumentRepository
}

func NewDocumentHandler(repo DocumentRepository) *DocumentHandler {
	return &DocumentHandler{repo: repo}
}

func (h *DocumentHandler) GetItems(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	join := c.DefaultQuery("join", "false") == "true"

	limitInt, err := strconv.Atoi(limit)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	results, err := h.repo.GetDocuments(limitInt, offsetInt, join)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *DocumentHandler) CreateItem(c *gin.Context) {
	document := &model.Document{}
	err := c.ShouldBindJSON(&document)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}
	document, err = h.repo.CreateDocument(document)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (h *DocumentHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	result, err := h.repo.GetDocument(idInt)
	if util.GinHandleError(c, err, http.StatusNotFound) {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DocumentHandler) UpdateItem(c *gin.Context) {
	document := model.Document{}
	err := c.ShouldBindJSON(&document)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	id := c.Param("id")

	document.ID, err = strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.UpdateDocument(&document)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DocumentHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.DeleteDocument(idInt)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
