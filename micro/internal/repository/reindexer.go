package repository

import (
	"fmt"
	"relorer/test_go_microservice/internal/model"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type ReindexerRepository struct {
	conn *reindexer.Reindexer
}

type ReindexerParams struct {
	Host     string
	Port     int
	Database string
}

func NewReindexerDB(params *ReindexerParams) (*ReindexerRepository, error) {
	conn := reindexer.NewReindex(fmt.Sprintf("cproto://%s:%d/%s", params.Host, params.Port, params.Database), reindexer.WithCreateDBIfMissing())

	if conn.Status().Err != nil {
		return nil, fmt.Errorf("failed to connect to Reindexer: %v", conn.Status().Err)
	}

	// Check connection to Reindexer
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Reindexer: %v", err)
	}

	db := &ReindexerRepository{
		conn: conn,
	}

	conn.OpenNamespace("documents", reindexer.DefaultNamespaceOptions(), model.Document{})

	return db, nil
}

func (r *ReindexerRepository) GetDocuments(limit, offset int) ([]*model.Document, error) {
	query := r.conn.Query("documents").Limit(limit).Offset(offset).Sort("id", true)
	data, err := query.Exec().FetchAll()
	if err != nil {
		return nil, err
	}

	documents := make([]*model.Document, len(data))
	for i, arg := range data {
		documents[i] = arg.(*model.Document)
	}

	return documents, nil
}

func (r *ReindexerRepository) CreateDocument(document *model.Document) (*model.Document, error) {
	_, err := r.conn.Insert("documents", document, "id=serial()")
	return document, err
}

func (r *ReindexerRepository) GetDocument(id int64) (*model.Document, error) {
	elem, found := r.conn.Query("documents").Where("id", reindexer.EQ, id).Get()

	if found {
		return elem.(*model.Document), nil
	}

	return nil, nil
}

func (r *ReindexerRepository) UpdateDocument(document *model.Document) error {
	_, err := r.conn.Update("documents", document)
	return err
}

func (r *ReindexerRepository) DeleteDocument(id int64) error {
	return r.conn.Delete("documents", &model.Document{ID: id})
}
