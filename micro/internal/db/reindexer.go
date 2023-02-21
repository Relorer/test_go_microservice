package db

import (
	"fmt"
	"relorer/test_go_microservice/internal/model"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type ReindexerDB struct {
	conn *reindexer.Reindexer
}

type ReindexerParams struct {
	Host     string
	Port     int
	Database string
}

func NewReindexerDB(params *ReindexerParams) (*ReindexerDB, error) {
	conn := reindexer.NewReindex(fmt.Sprintf("cproto://%s:%d/%s", params.Host, params.Port, params.Database), reindexer.WithCreateDBIfMissing())

	if conn.Status().Err != nil {
		return nil, fmt.Errorf("failed to connect to Reindexer: %v", conn.Status().Err)
	}

	// Check connection to Reindexer
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Reindexer: %v", err)
	}

	db := &ReindexerDB{
		conn: conn,
	}

	conn.OpenNamespace("documents", reindexer.DefaultNamespaceOptions(), model.Document{})

	return db, nil
}

func (r *ReindexerDB) GetDocuments(limit, offset int) ([]model.Document, error) {
	query := r.conn.Query("documents").Limit(limit).Offset(offset).Sort("id", true)
	documents, err := query.Exec().FetchAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(documents...)
	return nil, nil
}

func (r *ReindexerDB) CreateDocument(document *model.Document) (*model.Document, error) {
	// _, err := r.conn.Insert("documents").Values(document).Exec()
	// return err
	return nil, nil
}

func (r *ReindexerDB) GetDocument(id string) (*model.Document, error) {
	// query := r.conn.Query("documents").WhereString("id", reindexer.EQ, id)
	// var documents []model.Document
	// err := query.Exec().FetchAll(&documents)
	// if err != nil {
	// 	return nil, err
	// }
	// if len(documents) == 0 {
	// 	return nil, nil
	// }
	// return &documents[0], nil
	return nil, nil
}

func (r *ReindexerDB) UpdateDocument(document *model.Document) error {
	// _, err := r.conn.Update("documents").Set(document).WhereString("id", reindexer.EQ, document.ID).Exec()
	// return err
	return nil
}

func (r *ReindexerDB) DeleteDocument(id string) error {
	// _, err := r.conn.Delete("documents").WhereString("id", reindexer.EQ, id).Exec()
	// return err
	return nil
}
