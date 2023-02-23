package util

import (
	"relorer/test_go_microservice/internal/model"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateDocument(maxAuthorCount, totalAuthors int) *model.Document {
	ids := make([]int64, gofakeit.IntRange(1, maxAuthorCount))
	for i := range ids {
		ids[i] = int64(gofakeit.IntRange(1, totalAuthors))
	}

	return &model.Document{
		Title:      gofakeit.Sentence(5),
		Body:       gofakeit.Paragraph(5, 10, 5, "\n"),
		AuthorsIDs: ids,
	}
}

func GenerateAuthor(maxCommentsCount, totalComments int) *model.Author {

	ids := make([]int64, gofakeit.IntRange(1, maxCommentsCount))
	for i := range ids {
		ids[i] = int64(gofakeit.IntRange(1, totalComments))
	}

	return &model.Author{
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		Email:      gofakeit.Email(),
		Biography:  gofakeit.Paragraph(3, 5, 3, "\n"),
		Sort:       gofakeit.Int64(),
		CommentIDs: ids,
	}
}

func GenerateComment() *model.Comment {
	comment := &model.Comment{
		Text:     gofakeit.Paragraph(2, 4, 5, "\n"),
		Date:     gofakeit.Date(),
		Likes:    int64(gofakeit.IntRange(0, 5000)),
		Dislikes: int64(gofakeit.IntRange(0, 5000)),
		Sort:     gofakeit.Int64(),
	}
	return comment
}
