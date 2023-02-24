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

type MockDocumentRepo struct {
	mock.Mock
}

func (m *MockDocumentRepo) GetDocuments(limit, offset int, join bool) ([]*model.Document, error) {
	args := m.Called(limit, offset, join)
	return args.Get(0).([]*model.Document), args.Error(1)
}

func (m *MockDocumentRepo) CreateDocument(document *model.Document) (*model.Document, error) {
	args := m.Called(document)
	return args.Get(0).(*model.Document), args.Error(1)
}

func (m *MockDocumentRepo) GetDocument(id int64) (*model.Document, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Document), args.Error(1)
}

func (m *MockDocumentRepo) UpdateDocument(document *model.Document) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockDocumentRepo) DeleteDocument(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestDocumentHandler_GetItems(t *testing.T) {
	mockRepo := &MockDocumentRepo{}
	mockDoc := &model.Document{
		ID:    1,
		Title: "test title",
		Body:  "test content",
	}

	mockRepo.On("GetDocuments", 10, 0, false).Return([]*model.Document{mockDoc}, nil)

	handler := NewDocumentHandler(mockRepo)
	router := gin.Default()
	router.GET("/documents", handler.GetItems)

	req, err := http.NewRequest(http.MethodGet, "/documents", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var docs []*model.Document
	err = json.Unmarshal(rec.Body.Bytes(), &docs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(docs))
	assert.Equal(t, mockDoc, docs[0])
	mockRepo.AssertExpectations(t)
}

func TestDocumentHandler_GetItem(t *testing.T) {
	mockRepo := &MockDocumentRepo{}
	id := int64(1)
	mockDoc := &model.Document{
		ID:    id,
		Title: "test title",
		Body:  "test content",
	}

	mockRepo.On("GetDocument", id).Return(mockDoc, nil)

	handler := NewDocumentHandler(mockRepo)
	router := gin.Default()
	router.GET("/documents/:id", handler.GetItem)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/documents/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var doc *model.Document
	err = json.Unmarshal(rec.Body.Bytes(), &doc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, doc)
	mockRepo.AssertExpectations(t)
}

func TestDocumentHandler_CreateItem(t *testing.T) {
	mockRepo := &MockDocumentRepo{}
	mockDoc := &model.Document{
		ID:    1,
		Title: "test title",
		Body:  "test content",
	}

	mockRepo.On("CreateDocument", mock.AnythingOfType("*model.Document")).Return(mockDoc, nil)

	handler := NewDocumentHandler(mockRepo)
	router := gin.Default()
	router.POST("/documents", handler.CreateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(docJson))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdDoc model.Document
	err = json.Unmarshal(rec.Body.Bytes(), &createdDoc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockDoc, &createdDoc)
	mockRepo.AssertExpectations(t)
}

func TestDocumentHandler_UpdateItem(t *testing.T) {
	mockRepo := &MockDocumentRepo{}
	id := int64(1)
	mockDoc := &model.Document{
		ID:    id,
		Title: "Updated Document Title",
		Body:  "Updated Document Description",
	}

	mockRepo.On("UpdateDocument", mockDoc).Return(nil)

	handler := NewDocumentHandler(mockRepo)
	router := gin.Default()
	router.PUT("/documents/:id", handler.UpdateItem)

	docJson, err := json.Marshal(mockDoc)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/documents/%d", id), bytes.NewBuffer(docJson))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestDocumentHandler_DeleteItem(t *testing.T) {
	mockRepo := &MockDocumentRepo{}
	id := int64(1)

	handler := NewDocumentHandler(mockRepo)
	router := gin.Default()
	router.DELETE("/documents/:id", handler.DeleteItem)

	mockRepo.On("DeleteDocument", id).Return(nil)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/documents/%d", id), nil)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)

	mockRepo.AssertExpectations(t)
}
