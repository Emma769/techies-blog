package article

import (
	"context"
	"errors"
	"time"

	"github.com/emma769/techies-blog/internal/model"
	"github.com/emma769/techies-blog/internal/repository"
	"github.com/emma769/techies-blog/internal/services"
	"github.com/emma769/techies-blog/internal/utils/funclib"
)

type storer interface {
	CreateArticle(context.Context, model.ArticleParam) (model.Article, error)
	FindArticles(context.Context) ([]model.Article, error)
	FindArticle(context.Context, string) (model.Article, error)
	UpdateArticle(context.Context, model.Article) (model.Article, error)
	DeleteArticle(context.Context, string) error
}

type Service struct {
	store   storer
	timeout time.Duration
}

func NewService(store storer) *Service {
	return &Service{
		store:   store,
		timeout: 3 * time.Second,
	}
}

func (s *Service) Create(ctx context.Context, in model.ArticleIn) (model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	param := model.ArticleParam{
		Title:       in.Title,
		Slug:        funclib.Slugify(in.Title),
		Description: in.Description,
		Content:     in.Content,
	}

	article, err := s.store.CreateArticle(ctx, param)

	if err != nil && errors.Is(err, repository.ErrDuplicateKey) {
		return model.Article{}, services.ErrDuplicateKey
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (s *Service) FindAll(ctx context.Context) ([]model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.store.FindArticles(ctx)
}

func (s *Service) FindOne(ctx context.Context, slug string) (model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	article, err := s.store.FindArticle(ctx, slug)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return model.Article{}, services.ErrNotFound
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (s *Service) Update(
	ctx context.Context,
	stale model.Article,
	in model.ArticleIn,
) (model.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	stale.Title = in.Title
	stale.Slug = funclib.Slugify(in.Title)
	stale.Description = in.Description
	stale.Content = in.Content

	article, err := s.store.UpdateArticle(ctx, stale)

	if err != nil && errors.Is(err, repository.ErrUpdateConflict) {
		return model.Article{}, services.ErrUpdateConflict
	}

	if err != nil {
		return model.Article{}, err
	}

	return article, nil
}

func (s *Service) Delete(ctx context.Context, slug string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.store.DeleteArticle(ctx, slug)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return services.ErrNotFound
	}

	if err != nil {
		return err
	}

	return nil
}
