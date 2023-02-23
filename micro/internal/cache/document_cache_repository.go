package cache

import (
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/server/handler"
	"strconv"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type DocumentCacheRepository struct {
	c    Cache
	repo handler.DocumentRepository
}

func NewDocumentCacheRepository(c Cache, repo handler.DocumentRepository) *DocumentCacheRepository {
	return &DocumentCacheRepository{c: c, repo: repo}
}

func (r *DocumentCacheRepository) GetDocuments(limit, offset int, join bool) ([]*model.Document, error) {
	return r.repo.GetDocuments(limit, offset, join)
}

func (r *DocumentCacheRepository) CreateDocument(document *model.Document) (*model.Document, error) {
	return r.repo.CreateDocument(document)
}

func (r *DocumentCacheRepository) GetDocument(id int64) (*model.Document, error) {
	key := strconv.FormatInt(id, 10)
	cdoc, ok := r.c.Get(key)
	if ok {
		return cdoc.(*model.Document), nil
	}

	doc, err := r.repo.GetDocument(id)
	if err == nil {
		r.c.Set(key, doc)
		return doc, nil
	}

	return nil, nil
}

func (r *DocumentCacheRepository) UpdateDocument(document *model.Document) error {
	return r.repo.UpdateDocument(document)
}

func (r *DocumentCacheRepository) DeleteDocument(id int64) error {
	return r.repo.DeleteDocument(id)
}
