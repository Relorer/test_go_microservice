package util

import (
	"reflect"
	"relorer/test_go_microservice/internal/model"
	"testing"
	"time"
)

func TestToMaps(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "Document converting",
			args: args{
				obj: document,
			},
			want: documentWithMaps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMaps(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

var document = &model.Document{
	ID:         1,
	Title:      "My Document",
	Body:       "This is my document.",
	AuthorsIDs: []int64{1, 2},
	Authors: []*model.Author{
		{
			ID:         1,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "johndoe@example.com",
			Biography:  "John Doe is a writer and journalist.",
			CommentIDs: []int64{1, 2, 3},
			Sort:       1,
			Comments: []*model.Comment{
				{
					ID:       1,
					Text:     "Great article!",
					Date:     time.Date(2023, 2, 20, 8, 30, 0, 0, time.UTC),
					Likes:    10,
					Dislikes: 0,
					Sort:     1,
				},
				{
					ID:       2,
					Text:     "I really enjoyed reading this.",
					Date:     time.Date(2023, 2, 20, 9, 15, 0, 0, time.UTC),
					Likes:    5,
					Dislikes: 0,
					Sort:     2,
				},
				{
					ID:       3,
					Text:     "Thanks for sharing.",
					Date:     time.Date(2023, 2, 20, 10, 0, 0, 0, time.UTC),
					Likes:    3,
					Dislikes: 0,
					Sort:     3,
				},
			},
		},
		{
			ID:         2,
			FirstName:  "Jane",
			LastName:   "Doe",
			Email:      "janedoe@example.com",
			Biography:  "Jane Doe is a blogger and writer.",
			CommentIDs: []int64{4, 5, 6},
			Sort:       2,
			Comments: []*model.Comment{
				{
					ID:       4,
					Text:     "Interesting insights.",
					Date:     time.Date(2023, 2, 21, 11, 30, 0, 0, time.UTC),
					Likes:    8,
					Dislikes: 0,
					Sort:     1,
				},
				{
					ID:       5,
					Text:     "I disagree with some of your points.",
					Date:     time.Date(2023, 2, 21, 12, 15, 0, 0, time.UTC),
					Likes:    2,
					Dislikes: 5,
					Sort:     2,
				},
				{
					ID:       6,
					Text:     "This was helpful, thank you.",
					Date:     time.Date(2023, 2, 21, 13, 0, 0, 0, time.UTC),
					Likes:    5,
					Dislikes: 0,
					Sort:     3,
				},
			},
		},
	},
}

var documentWithMaps = map[string]interface{}{
	"id":          int64(1),
	"title":       "My Document",
	"body":        "This is my document.",
	"authors_ids": []interface{}{int64(1), int64(2)},
	"authors": []interface{}{
		map[string]interface{}{
			"id":           int64(1),
			"first_name":   "John",
			"last_name":    "Doe",
			"email":        "johndoe@example.com",
			"biography":    "John Doe is a writer and journalist.",
			"comments_ids": []interface{}{int64(1), int64(2), int64(3)},
			"sort":         int64(1),
			"comments": []interface{}{
				map[string]interface{}{
					"id":       int64(1),
					"text":     "Great article!",
					"date":     time.Date(2023, 2, 20, 8, 30, 0, 0, time.UTC),
					"likes":    int64(10),
					"dislikes": int64(0),
					"sort":     int64(1),
				},
				map[string]interface{}{
					"id":       int64(2),
					"text":     "I really enjoyed reading this.",
					"date":     time.Date(2023, 2, 20, 9, 15, 0, 0, time.UTC),
					"likes":    int64(5),
					"dislikes": int64(0),
					"sort":     int64(2),
				},
				map[string]interface{}{
					"id":       int64(3),
					"text":     "Thanks for sharing.",
					"date":     time.Date(2023, 2, 20, 10, 0, 0, 0, time.UTC),
					"likes":    int64(3),
					"dislikes": int64(0),
					"sort":     int64(3),
				},
			},
		},
		map[string]interface{}{
			"id":           int64(2),
			"first_name":   "Jane",
			"last_name":    "Doe",
			"email":        "janedoe@example.com",
			"biography":    "Jane Doe is a blogger and writer.",
			"comments_ids": []interface{}{int64(4), int64(5), int64(6)},
			"sort":         int64(2),
			"comments": []interface{}{
				map[string]interface{}{
					"id":       int64(4),
					"text":     "Interesting insights.",
					"date":     time.Date(2023, 2, 21, 11, 30, 0, 0, time.UTC),
					"likes":    int64(8),
					"dislikes": int64(0),
					"sort":     int64(1),
				},
				map[string]interface{}{
					"id":       int64(5),
					"text":     "I disagree with some of your points.",
					"date":     time.Date(2023, 2, 21, 12, 15, 0, 0, time.UTC),
					"likes":    int64(2),
					"dislikes": int64(5),
					"sort":     int64(2),
				},
				map[string]interface{}{
					"id":       int64(6),
					"text":     "This was helpful, thank you.",
					"date":     time.Date(2023, 2, 21, 13, 0, 0, 0, time.UTC),
					"likes":    int64(5),
					"dislikes": int64(0),
					"sort":     int64(3),
				},
			},
		},
	},
}
