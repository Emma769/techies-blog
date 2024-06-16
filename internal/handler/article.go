package handler

import (
	"errors"
	"net/http"

	"github.com/emma769/techies-blog/internal/model"
	"github.com/emma769/techies-blog/internal/services"
	"github.com/emma769/techies-blog/internal/utils/funclib"
)

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) error {
	in, err := BIND[model.ArticleIn](h, r)
	if err != nil {
		return NewHandlerError(422, err.Error())
	}

	article, err := h.Article.Create(r.Context(), in)

	if err != nil && errors.Is(err, services.ErrDuplicateKey) {
		return NewHandlerError(400, "slug already exists")
	}

	if err != nil {
		return err
	}

	return h.JSON(w, 201, model.CreateArticleOut(article))
}

func (h *Handler) FindArticles(w http.ResponseWriter, r *http.Request) error {
	articles, err := h.Article.FindAll(r.Context())
	if err != nil {
		return err
	}

	return h.JSON(w, 200, funclib.Map(articles, model.CreateArticleOut))
}

func (h *Handler) FindArticle(w http.ResponseWriter, r *http.Request) error {
	slug := h.Param(r, "slug")

	article, err := h.Article.FindOne(r.Context(), slug)

	if err != nil && errors.Is(err, services.ErrNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return err
	}

	return h.JSON(w, 200, model.CreateArticleOut(article))
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) error {
	slug := h.Param(r, "slug")

	stale, err := h.Article.FindOne(r.Context(), slug)

	if err != nil && errors.Is(err, services.ErrNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return err
	}

	in, err := BIND[model.ArticleIn](h, r)
	if err != nil {
		return NewHandlerError(422, err.Error())
	}

	article, err := h.Article.Update(r.Context(), stale, in)

	if err != nil && errors.Is(err, services.ErrUpdateConflict) {
		return ErrConflict
	}

	if err != nil {
		return err
	}

	return h.JSON(w, 200, model.CreateArticleOut(article))
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) error {
	slug := h.Param(r, "slug")

	err := h.Article.Delete(r.Context(), slug)

	if err != nil && errors.Is(err, services.ErrNotFound) {
		return ErrNotFound
	}

	if err != nil {
		return err
	}

	return h.Status(w, 204)
}
