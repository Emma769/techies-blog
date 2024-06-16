package model

import (
	"time"

	"github.com/emma769/techies-blog/internal/utils/funclib"
	"github.com/emma769/techies-blog/internal/utils/validator"
)

type Article struct {
	ArticleID   string     `json:"article_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	Version     int64      `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type ArticleParam struct {
	Title       string
	Slug        string
	Description string
	Content     string
}

type ArticleIn struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

var (
	gte2  = funclib.GteCUR(2)
	lte50 = funclib.LteCUR(150)
)

func (in ArticleIn) Validate() error {
	return validator.Check(
		validator.New(),
		in,
		func(t ArticleIn) (string, bool) {
			return "title:cannot be blank", funclib.NonWhiteSpace(t.Title)
		},
		func(t ArticleIn) (string, bool) {
			return "title:must contain at least 2 characters", gte2(len(t.Title))
		},
		func(t ArticleIn) (string, bool) {
			return "title:must contain at most 150 characters", lte50(len(t.Title))
		},
		func(t ArticleIn) (string, bool) {
			return "content:cannot be blank", funclib.NonWhiteSpace(t.Content)
		},
	)
}

type ArticleOut struct {
	ArticleID   string     `json:"article_id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	HTML        string     `json:"html"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func CreateArticleOut(article Article) ArticleOut {
	sanitized := funclib.SanitizeHTML(funclib.ParseToHTML([]byte(article.Content)))

	return ArticleOut{
		ArticleID:   article.ArticleID,
		Title:       article.Title,
		Slug:        article.Slug,
		Description: article.Description,
		HTML:        sanitized,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}
