package pgstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/emma769/techies-blog/internal/model"
	"github.com/emma769/techies-blog/internal/repository"
)

func (q *Queries) CreateArticle(
	ctx context.Context,
	param model.ArticleParam,
) (model.Article, error) {
	query := `INSERT INTO articles (title, slug, description, content)
  VALUES ($1, $2, $3, $4) RETURNING article_id, title, slug, description, 
  content, version, created_at, updated_at;`

	row := q.db.QueryRowContext(
		ctx,
		query,
		param.Title,
		param.Slug,
		param.Description,
		param.Content,
	)

	var article model.Article

	err := scanArticleRow(row, &article)

	if err != nil && repository.DuplKey(err.Error()) {
		return model.Article{}, repository.ErrDuplicateKey
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (q *Queries) FindArticles(ctx context.Context) ([]model.Article, error) {
	query := `SELECT article_id, title, slug, description, content, version,
  created_at, updated_at FROM articles;`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}

	for rows.Next() {
		var article model.Article
		err := scanArticleRow(rows, &article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (q *Queries) FindArticle(ctx context.Context, slug string) (model.Article, error) {
	query := `SELECT article_id, title, slug, description, content, version,
  created_at, updated_at FROM articles WHERE slug = $1;`

	row := q.db.QueryRowContext(ctx, query, slug)

	var article model.Article

	err := scanArticleRow(row, &article)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.Article{}, repository.ErrNotFound
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (q *Queries) UpdateArticle(ctx context.Context, updated model.Article) (model.Article, error) {
	query := `UPDATE articles SET title = $1, slug = $2, description = $3, content = $4
  WHERE article_id = $5 AND version = $6 RETURNING article_id, title, 
  slug, description, content, version, created_at, updated_at;`

	row := q.db.QueryRowContext(
		ctx,
		query,
		updated.Title,
		updated.Slug,
		updated.Description,
		updated.Content,
		updated.ArticleID,
		updated.Version,
	)

	var article model.Article

	err := scanArticleRow(row, &article)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return model.Article{}, repository.ErrUpdateConflict
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (q *Queries) DeleteArticle(ctx context.Context, slug string) error {
	query := `DELETE FROM articles WHERE slug = $1;`

	result, err := q.db.ExecContext(ctx, query, slug)
	if err != nil {
		return err
	}

	aff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return repository.ErrNotFound
	}

	return nil
}

type scanner interface {
	Scan(...any) error
}

func scanArticleRow[S scanner](s S, article *model.Article) error {
	return s.Scan(
		&article.ArticleID,
		&article.Title,
		&article.Slug,
		&article.Description,
		&article.Content,
		&article.Version,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
}
