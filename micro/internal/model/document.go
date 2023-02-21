package model

type Document struct {
	ID    int64  `json:"id" reindex:"id,,pk"`
	Title string `json:"title" reindex:"title,fuzzytext"`
	Body  string `json:"body" reindex:"body,text"`
}
