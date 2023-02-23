package handler

import (
	"net/http"
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentRepository interface {
	GetComments(limit, offset int, join bool) ([]*model.Comment, error)
	CreateComment(comment *model.Comment) (*model.Comment, error)
	GetComment(id int64) (*model.Comment, error)
	UpdateComment(comment *model.Comment) error
	DeleteComment(id int64) error
}

type CommentHandler struct {
	repo CommentRepository
}

func NewCommentHandler(repo CommentRepository) *CommentHandler {
	return &CommentHandler{repo: repo}
}

func (h *CommentHandler) GetItems(c *gin.Context) {
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

	results, err := h.repo.GetComments(limitInt, offsetInt, join)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *CommentHandler) CreateItem(c *gin.Context) {
	comment := &model.Comment{}
	err := c.ShouldBindJSON(&comment)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}
	comment, err = h.repo.CreateComment(comment)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	result, err := h.repo.GetComment(idInt)
	if util.GinHandleError(c, err, http.StatusNotFound) {
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *CommentHandler) UpdateItem(c *gin.Context) {
	comment := model.Comment{}
	err := c.ShouldBindJSON(&comment)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	id := c.Param("id")

	comment.ID, err = strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.UpdateComment(&comment)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CommentHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if util.GinHandleError(c, err, http.StatusBadRequest) {
		return
	}

	err = h.repo.DeleteComment(idInt)
	if util.GinHandleError(c, err, http.StatusInternalServerError) {
		return
	}

	c.Status(http.StatusNoContent)
}
