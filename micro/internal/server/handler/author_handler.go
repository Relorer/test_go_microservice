package handler

import (
	"net/http"
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthorRepository interface {
	GetAuthors(limit, offset int, join bool) ([]*model.Author, error)
	CreateAuthor(author *model.Author) (*model.Author, error)
	GetAuthor(id int64) (*model.Author, error)
	UpdateAuthor(author *model.Author) error
	DeleteAuthor(id int64) error
}

type AuthorHandler struct {
	repo AuthorRepository
}

func NewAuthorHandler(repo AuthorRepository) *AuthorHandler {
	return &AuthorHandler{repo: repo}
}

func (h *AuthorHandler) GetItems(c *gin.Context) {
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

	results, err := h.repo.GetAuthors(limitInt, offsetInt, join)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *AuthorHandler) CreateItem(c *gin.Context) {
	author := &model.Author{}
	err := c.ShouldBindJSON(&author)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}
	author, err = h.repo.CreateAuthor(author)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusCreated, author)
}

func (h *AuthorHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	result, err := h.repo.GetAuthor(idInt)
	if util.GinHandleError(c, err, http.StatusNotFound) {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthorHandler) UpdateItem(c *gin.Context) {
	author := model.Author{}
	err := c.ShouldBindJSON(&author)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	id := c.Param("id")

	author.ID, err = strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.UpdateAuthor(&author)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthorHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.DeleteAuthor(idInt)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
