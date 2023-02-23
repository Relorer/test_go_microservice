package model

import "time"

type Document struct {
	ID         int64     `json:"id" reindex:"id,,pk"`
	Title      string    `json:"title" reindex:"title,fuzzytext"`
	Body       string    `json:"body" reindex:"body,text"`
	AuthorsIDs []int64   `json:"authors_ids" reindex:"authors_ids"`
	Authors    []*Author `json:"authors" reindex:"authors,,joined"`
}

type Author struct {
	ID         int64      `json:"id" reindex:"id,,pk"`
	FirstName  string     `json:"first_name" reindex:"first_name,text"`
	LastName   string     `json:"last_name" reindex:"last_name,text"`
	Email      string     `json:"email" reindex:"email,text"`
	Sort       int64      `json:"sort" reindex:"sort"`
	Biography  string     `json:"biography" reindex:"biography,text"`
	CommentIDs []int64    `json:"comments_ids" reindex:"comments_ids"`
	Comments   []*Comment `json:"comments" reindex:"comments,,joined"`
}

type Comment struct {
	ID       int64     `json:"id" reindex:"id,,pk"`
	Text     string    `json:"text" reindex:"text,text"`
	Date     time.Time `json:"date" reindex:"date,tree"`
	Likes    int64     `json:"likes" reindex:"likes"`
	Dislikes int64     `json:"dislikes" reindex:"dislikes"`
	Sort     int64     `json:"sort" reindex:"sort"`
}
