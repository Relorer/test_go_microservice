package database

import (
	"fmt"
	"log"
	"relorer/test_go_microservice/internal/model"
	"relorer/test_go_microservice/internal/util"
	"time"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

const DocumentsNamespace string = "documents"
const AuthorsNamespace string = "authors"
const CommentsNamespace string = "comments"

type ReindexerParams struct {
	Host             string
	Port             int
	Database         string
	GenerateTestData bool
}

func ReindexerConnectWithRetry(params *ReindexerParams, delay time.Duration) *reindexer.Reindexer {
	for {
		db, err := NewReindexerRepository(params)

		if err != nil {
			log.Printf("Error connecting: %s. Retry in %s", err, delay.String())
		} else {
			return db
		}

		time.Sleep(delay)
	}
}

func NewReindexerRepository(params *ReindexerParams) (*reindexer.Reindexer, error) {
	conn := reindexer.NewReindex(fmt.Sprintf("cproto://%s:%d/%s", params.Host, params.Port, params.Database), reindexer.WithCreateDBIfMissing())

	if conn.Status().Err != nil {
		return nil, fmt.Errorf("failed to connect to Reindexer: %v", conn.Status().Err)
	}

	// Check connection to Reindexer
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Reindexer: %v", err)
	}

	conn.OpenNamespace(DocumentsNamespace, reindexer.DefaultNamespaceOptions(), model.Document{})
	conn.OpenNamespace(AuthorsNamespace, reindexer.DefaultNamespaceOptions(), model.Author{})
	conn.OpenNamespace(CommentsNamespace, reindexer.DefaultNamespaceOptions(), model.Comment{})

	if params.GenerateTestData {
		documentsCount := conn.Query(DocumentsNamespace).Exec().Count()
		authorsCount := conn.Query(AuthorsNamespace).Exec().Count()
		commentsCount := conn.Query(CommentsNamespace).Exec().Count()
		if commentsCount+authorsCount+documentsCount == 0 {
			generateTestData(conn)
		}
	}

	return conn, nil
}

func generateTestData(conn *reindexer.Reindexer) {
	documentsCount := 100000
	authorsCount := 10000
	commentsCount := 10000

	for i := 0; i < commentsCount; i++ {
		comment := util.GenerateComment()
		comment.ID = int64(i)
		_, err := conn.Insert(CommentsNamespace, comment)
		if err != nil {
			log.Printf("Error adding comment: %s", err)
		}
	}

	for i := 0; i < authorsCount; i++ {
		author := util.GenerateAuthor(10, commentsCount)
		author.ID = int64(i)
		_, err := conn.Insert(AuthorsNamespace, author)
		if err != nil {
			log.Printf("Error adding author: %s", err)
		}
	}

	for i := 0; i < documentsCount; i++ {
		document := util.GenerateDocument(10, authorsCount)
		document.ID = int64(i)
		_, err := conn.Insert(DocumentsNamespace, document)
		if err != nil {
			log.Printf("Error adding document: %s", err)
		}
	}
}
