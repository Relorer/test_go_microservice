package database

import (
	"relorer/test_go_microservice/internal/model"

	"github.com/restream/reindexer/v3"
	_ "github.com/restream/reindexer/v3/bindings/cproto"
)

type CommentReindexerRepository struct {
	conn *reindexer.Reindexer
}

func NewCommentReindexerRepository(reindexer *reindexer.Reindexer) *CommentReindexerRepository {
	return &CommentReindexerRepository{conn: reindexer}
}

func (r *CommentReindexerRepository) GetComments(limit, offset int, join bool) ([]*model.Comment, error) {
	query := r.conn.Query(CommentsNamespace).Limit(limit).Offset(offset).Sort("id", true)
	data, err := query.Exec().FetchAll()
	if err != nil {
		return nil, err
	}

	comments := make([]*model.Comment, len(data))
	for i, arg := range data {
		comments[i] = arg.(*model.Comment)
	}

	return comments, nil
}

func (r *CommentReindexerRepository) CreateComment(comment *model.Comment) (*model.Comment, error) {
	_, err := r.conn.Insert(CommentsNamespace, comment, "id=serial()")
	return comment, err
}

func (r *CommentReindexerRepository) GetComment(id int64) (*model.Comment, error) {
	elem, found := r.conn.Query(CommentsNamespace).Where("id", reindexer.EQ, id).Get()

	if !found {
		return nil, nil
	}

	return elem.(*model.Comment), nil
}

func (r *CommentReindexerRepository) UpdateComment(comment *model.Comment) error {
	_, err := r.conn.Update(CommentsNamespace, comment)
	return err
}

func (r *CommentReindexerRepository) DeleteComment(id int64) error {
	return r.conn.Delete(CommentsNamespace, &model.Comment{ID: id})
}
