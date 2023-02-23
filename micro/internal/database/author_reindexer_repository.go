package database

import (
	"relorer/test_go_microservice/internal/model"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type AuthorReindexerRepository struct {
	conn *reindexer.Reindexer
}

func NewAuthorReindexerRepository(reindexer *reindexer.Reindexer) *AuthorReindexerRepository {
	return &AuthorReindexerRepository{conn: reindexer}
}

func (r *AuthorReindexerRepository) GetAuthors(limit, offset int, join bool) ([]*model.Author, error) {
	query := r.conn.Query(AuthorsNamespace).Limit(limit).Offset(offset).Sort("sort", true)

	if join {
		query = query.Join(r.conn.Query(CommentsNamespace).Sort("sort", true), CommentsNamespace).
			On("comments_ids", reindexer.SET, "id")
	}

	data, err := query.Exec().FetchAll()

	if err != nil {
		return nil, err
	}

	authors := make([]*model.Author, len(data))
	for i, arg := range data {
		authors[i] = arg.(*model.Author)
	}

	return authors, nil
}

func (r *AuthorReindexerRepository) CreateAuthor(author *model.Author) (*model.Author, error) {
	_, err := r.conn.Insert(AuthorsNamespace, author, "id=serial()")
	return author, err
}

func (r *AuthorReindexerRepository) GetAuthor(id int64) (*model.Author, error) {
	elem, err := r.conn.Query(AuthorsNamespace).Where("id", reindexer.EQ, id).
		Join(r.conn.Query(CommentsNamespace).Sort("sort", true), CommentsNamespace).
		On("comments_ids", reindexer.SET, "id").Exec().FetchOne()

	if err != nil {
		return nil, err
	}

	return elem.(*model.Author), nil
}

func (r *AuthorReindexerRepository) UpdateAuthor(author *model.Author) error {
	_, err := r.conn.Update(AuthorsNamespace, author)
	return err
}

func (r *AuthorReindexerRepository) DeleteAuthor(id int64) error {
	return r.conn.Delete(AuthorsNamespace, &model.Author{ID: id})
}
