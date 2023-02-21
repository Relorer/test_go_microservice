package model

type Document struct {
	ID    string `json:"id" reindex:"id,,pk"`
	Title string `json:"title" reindex:"title,fuzzytext"`
	Body  string `json:"body" reindex:"body,text"`
}
