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

type MockAuthorRepo struct {
	mock.Mock
}

func (m *MockAuthorRepo) GetAuthors(limit, offset int, join bool) ([]*model.Author, error) {
	args := m.Called(limit, offset, join)
	return args.Get(0).([]*model.Author), args.Error(1)
}

func (m *MockAuthorRepo) CreateAuthor(author *model.Author) (*model.Author, error) {
	args := m.Called(author)
	return args.Get(0).(*model.Author), args.Error(1)
}

func (m *MockAuthorRepo) GetAuthor(id int64) (*model.Author, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Author), args.Error(1)
}

func (m *MockAuthorRepo) UpdateAuthor(author *model.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

func (m *MockAuthorRepo) DeleteAuthor(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestAuthorHandler_GetItems(t *testing.T) {
	mockRepo := &MockAuthorRepo{}
	mockDoc := &model.Author{
		ID:        1,
		FirstName: "First name",
		LastName:  "Last name",
	}

	mockRepo.On("GetAuthors", 10, 0, false).Return([]*model.Author{mockDoc}, nil)

	handler := NewAuthorHandler(mockRepo)
	router := gin.Default()
	router.GET("/authors", handler.GetItems)

	req, err := http.NewRequest(http.MethodGet, "/authors", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var docs []*model.Author
	err = json.Unmarshal(rec.Body.Bytes(), &docs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(docs))
	assert.Equal(t, mockDoc, docs[0])
	mockRepo.AssertExpectations(t)
}

func TestAuthorHandler_GetItem(t *testing.T) {
	mockRepo := &MockAuthorRepo{}
	id := int64(1)
	mockDoc := &model.Author{
		ID:        id,
		FirstName: "First name",
		LastName:  "Last name",
	}

	mockRepo.On("GetAuthor", id).Return(mockDoc, nil)

	handler := NewAuthorHandler(mockRepo)
	router := gin.Default()
	router.GET("/authors/:id", handler.GetItem)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/authors/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var doc *model.Author
	err = json.Unmarshal(rec.Body.Bytes(), &doc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, doc)
	mockRepo.AssertExpectations(t)
}

func TestAuthorHandler_CreateItem(t *testing.T) {
	mockRepo := &MockAuthorRepo{}
	mockDoc := &model.Author{
		ID:        1,
		FirstName: "First name",
		LastName:  "Last name",
	}

	mockRepo.On("CreateAuthor", mock.AnythingOfType("*model.Author")).Return(mockDoc, nil)

	handler := NewAuthorHandler(mockRepo)
	router := gin.Default()
	router.POST("/authors", handler.CreateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/authors", bytes.NewBuffer(docJson))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdDoc model.Author
	err = json.Unmarshal(rec.Body.Bytes(), &createdDoc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, &createdDoc)
	mockRepo.AssertExpectations(t)
}

func TestAuthorHandler_UpdateItem(t *testing.T) {
	mockRepo := &MockAuthorRepo{}
	id := int64(1)
	mockDoc := &model.Author{
		ID:        id,
		FirstName: "First name",
		LastName:  "Last name",
	}

	mockRepo.On("UpdateAuthor", mockDoc).Return(nil)

	handler := NewAuthorHandler(mockRepo)
	router := gin.Default()
	router.PUT("/authors/:id", handler.UpdateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/authors/%d", id), bytes.NewBuffer(docJson))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestAuthorHandler_DeleteItem(t *testing.T) {
	mockRepo := &MockAuthorRepo{}
	id := int64(1)

	handler := NewAuthorHandler(mockRepo)
	router := gin.Default()
	router.DELETE("/authors/:id", handler.DeleteItem)

	mockRepo.On("DeleteAuthor", id).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/authors/%d", id), nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)

	mockRepo.AssertExpectations(t)
}
