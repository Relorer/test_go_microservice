package util

import (
	"relorer/test_go_microservice/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateDocument(t *testing.T) {
	type args struct {
		maxAuthorCount int
		totalAuthors   int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "generate document 2 5",
			args: args{
				maxAuthorCount: 2,
				totalAuthors:   5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateDocument(tt.args.maxAuthorCount, tt.args.totalAuthors)
			assert.NotNil(t, got)
			assert.IsType(t, &model.Document{}, got)

			assert.LessOrEqual(t, len(got.AuthorsIDs), 2)

			for _, id := range got.AuthorsIDs {
				assert.Less(t, id, int64(tt.args.totalAuthors))
			}
		})
	}
}

func TestGenerateAuthor(t *testing.T) {
	type args struct {
		maxCommentsCount int
		totalComments    int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "generate author 2 5",
			args: args{
				maxCommentsCount: 2,
				totalComments:    5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAuthor(tt.args.maxCommentsCount, tt.args.totalComments)
			assert.NotNil(t, got)
			assert.IsType(t, &model.Author{}, got)

			assert.LessOrEqual(t, len(got.CommentIDs), 2)

			for _, id := range got.CommentIDs {
				assert.Less(t, id, int64(tt.args.totalComments))
			}
		})
	}
}

func TestGenerateComment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "generate comment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateComment()
			assert.NotNil(t, got)
			assert.IsType(t, &model.Comment{}, got)
		})
	}
}
