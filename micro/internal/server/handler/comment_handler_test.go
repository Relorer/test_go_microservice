package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"relorer/test_go_microservice/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) GetComments(limit, offset int, join bool) ([]*model.Comment, error) {
	args := m.Called(limit, offset, join)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) CreateComment(comment *model.Comment) (*model.Comment, error) {
	args := m.Called(comment)
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) GetComment(id int64) (*model.Comment, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockCommentRepo) UpdateComment(comment *model.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockCommentRepo) DeleteComment(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCommentHandler_GetItems(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockDoc := &model.Comment{
		ID:    1,
		Text:  "test title",
		Likes: int64(123),
	}

	mockRepo.On("GetComments", 10, 0, false).Return([]*model.Comment{mockDoc}, nil)

	handler := NewCommentHandler(mockRepo)
	router := gin.Default()
	router.GET("/comments", handler.GetItems)

	req, err := http.NewRequest(http.MethodGet, "/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var docs []*model.Comment
	err = json.Unmarshal(rec.Body.Bytes(), &docs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(docs))
	assert.Equal(t, mockDoc, docs[0])
	mockRepo.AssertExpectations(t)
}

func TestCommentHandler_GetItem(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	id := int64(1)
	mockDoc := &model.Comment{
		ID:    id,
		Text:  "test title",
		Likes: int64(123),
	}

	mockRepo.On("GetComment", id).Return(mockDoc, nil)

	handler := NewCommentHandler(mockRepo)
	router := gin.Default()
	router.GET("/comments/:id", handler.GetItem)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/comments/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var doc *model.Comment
	err = json.Unmarshal(rec.Body.Bytes(), &doc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, doc)
	mockRepo.AssertExpectations(t)
}

func TestCommentHandler_CreateItem(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	mockDoc := &model.Comment{
		ID:    1,
		Text:  "test title",
		Likes: int64(123),
	}

	mockRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(mockDoc, nil)

	handler := NewCommentHandler(mockRepo)
	router := gin.Default()
	router.POST("/comments", handler.CreateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(docJson))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdDoc model.Comment
	err = json.Unmarshal(rec.Body.Bytes(), &createdDoc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, &createdDoc)
	mockRepo.AssertExpectations(t)
}

func TestCommentHandler_UpdateItem(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	id := int64(1)
	mockDoc := &model.Comment{
		ID:    id,
		Text:  "test title",
		Likes: int64(123),
	}

	mockRepo.On("UpdateComment", mockDoc).Return(nil)

	handler := NewCommentHandler(mockRepo)
	router := gin.Default()
	router.PUT("/comments/:id", handler.UpdateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/comments/%d", id), bytes.NewBuffer(docJson))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestCommentHandler_DeleteItem(t *testing.T) {
	mockRepo := &MockCommentRepo{}
	id := int64(1)

	handler := NewCommentHandler(mockRepo)
	router := gin.Default()
	router.DELETE("/comments/:id", handler.DeleteItem)

	mockRepo.On("DeleteComment", id).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/comments/%d", id), nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)

	mockRepo.AssertExpectations(t)
}
