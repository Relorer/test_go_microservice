package util

import (
	"relorer/test_go_microservice/internal/model"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateDocument() *model.Document {
	return &model.Document{
		Title: gofakeit.Sentence(5),
		Body:  gofakeit.Paragraph(5, 10, 5, "\n"),
	}
}
