package repository

import (
	"log"
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/server"
	"strconv"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}

type CacheRepository struct {
	c    Cache
	repo server.Repository
}

func NewCacheRepository(c Cache, repo server.Repository) *CacheRepository {
	return &CacheRepository{c: c, repo: repo}
}

func (r *CacheRepository) GetDocuments(limit, offset int) ([]*model.Document, error) {
	return r.repo.GetDocuments(limit, offset)
}

func (r *CacheRepository) CreateDocument(document *model.Document) (*model.Document, error) {
	return r.repo.CreateDocument(document)
}

func (r *CacheRepository) GetDocument(id int64) (*model.Document, error) {
	log.Println("start")
	key := strconv.FormatInt(id, 10)
	log.Println(key)
	cdoc, ok := r.c.Get(key)
	if ok {
		return cdoc.(*model.Document), nil
	}
	log.Println(key)

	doc, err := r.repo.GetDocument(id)
	if err == nil {
		r.c.Set(key, doc)
		return doc, nil
	}

	return nil, nil
}

func (r *CacheRepository) UpdateDocument(document *model.Document) error {
	return r.repo.UpdateDocument(document)
}

func (r *CacheRepository) DeleteDocument(id int64) error {
	return r.repo.DeleteDocument(id)
}
