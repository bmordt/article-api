package models

import "time"

type Article struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Body  string    `json:"body"`
	Tags  []string  `json:"tags"`
}

type CreateArticleReq struct {
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type ArticleResp struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}
