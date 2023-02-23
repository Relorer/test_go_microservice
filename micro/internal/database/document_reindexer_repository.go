package database

import (
	"relorer/test_go_microservice/internal/model"
	"sync"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type DocumentReindexerRepository struct {
	conn *reindexer.Reindexer
}

func NewDocumentReindexerRepository(reindexer *reindexer.Reindexer) *DocumentReindexerRepository {
	return &DocumentReindexerRepository{conn: reindexer}
}

func (r *DocumentReindexerRepository) GetDocuments(limit, offset int, join bool) ([]*model.Document, error) {
	query := r.conn.Query(DocumentsNamespace).Limit(limit).Offset(offset).Sort("title", true)
	data, err := query.Exec().FetchAll()
	if err != nil {
		return nil, err
	}

	documents := make([]*model.Document, len(data))
	for i, arg := range data {
		documents[i] = arg.(*model.Document)
	}

	if join {
		r.loadNestedFields(documents)
	}

	return documents, nil
}

func (r *DocumentReindexerRepository) loadNestedFields(documents []*model.Document) {
	var wg sync.WaitGroup
	wg.Add(len(documents))
	for _, doc := range documents {
		func(doc *model.Document) {
			defer wg.Done()

			items, err := r.conn.Query(AuthorsNamespace).WhereInt64("id", reindexer.EQ, doc.AuthorsIDs...).Sort("sort", true).
				Join(r.conn.Query(CommentsNamespace).Sort("sort", true), CommentsNamespace).
				On("comments_ids", reindexer.SET, "id").Exec().FetchAll()

			if err != nil {
				return
			}

			doc.Authors = make([]*model.Author, len(items))
			for i, value := range items {
				doc.Authors[i] = value.(*model.Author)
			}
		}(doc)
	}
	wg.Wait()
}

func (r *DocumentReindexerRepository) CreateDocument(document *model.Document) (*model.Document, error) {
	_, err := r.conn.Insert(DocumentsNamespace, document, "id=serial()")
	return document, err
}

func (r *DocumentReindexerRepository) GetDocument(id int64) (*model.Document, error) {
	elem, found := r.conn.Query(DocumentsNamespace).Where("id", reindexer.EQ, id).Get()

	if !found {
		return nil, nil
	}

	doc := elem.(*model.Document)

	items, err := r.conn.Query(AuthorsNamespace).WhereInt64("id", reindexer.EQ, doc.AuthorsIDs...).Sort("sort", true).
		Join(r.conn.Query(CommentsNamespace).Sort("sort", true), CommentsNamespace).
		On("comments_ids", reindexer.SET, "id").Exec().FetchAll()

	if err != nil {
		return nil, err
	}

	doc.Authors = make([]*model.Author, len(items))
	for i, value := range items {
		doc.Authors[i] = value.(*model.Author)
	}

	return doc, nil
}

func (r *DocumentReindexerRepository) UpdateDocument(document *model.Document) error {
	_, err := r.conn.Update(DocumentsNamespace, document)
	return err
}

func (r *DocumentReindexerRepository) DeleteDocument(id int64) error {
	return r.conn.Delete(DocumentsNamespace, &model.Document{ID: id})
}
