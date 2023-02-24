package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConfigureCRUDRoutes(t *testing.T) {

	router := gin.New()
	ConfigureCRUDRoutes(router, CRUDGroup{&mockHandler{}, "items"})

	tests := []struct {
		method string
		url    string
		status int
		body   string
	}{
		{"GET", "/items", http.StatusOK, "{\"message\":\"get items\"}"},
		{"GET", "/items/123", http.StatusOK, "{\"message\":\"get item\"}"},
		{"POST", "/items", http.StatusOK, "{\"message\":\"create item\"}"},
		{"PUT", "/items/123", http.StatusOK, "{\"message\":\"update item\"}"},
		{"DELETE", "/items/123", http.StatusOK, "{\"message\":\"delete item\"}"},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.url, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, tt.status, resp.Code)
		assert.Equal(t, tt.body, resp.Body.String())
	}

}

type mockHandler struct{}

func (h *mockHandler) GetItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get items"})
}

func (h *mockHandler) GetItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get item"})
}

func (h *mockHandler) CreateItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create item"})
}

func (h *mockHandler) UpdateItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update item"})
}

func (h *mockHandler) DeleteItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete item"})
}
