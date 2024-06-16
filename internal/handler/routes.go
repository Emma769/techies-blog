package handler

import "github.com/go-chi/chi/v5"

func (h *Handler) Register(mux *chi.Mux) {
	mux.Get("/api/healthz", wrap(h.Healthz))

	mux.Post("/api/articles", wrap(h.CreateArticle))
	mux.Get("/api/articles", wrap(h.FindArticles))
	mux.Get("/api/articles/{slug}", wrap(h.FindArticle))
	mux.Put("/api/articles/{slug}", wrap(h.UpdateArticle))
	mux.Delete("/api/articles/{slug}", wrap(h.DeleteArticle))
}
