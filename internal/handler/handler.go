package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/emma769/techies-blog/internal/services/article"
	"github.com/emma769/techies-blog/internal/utils/funclib"
)

type Services struct {
	Article *article.Service
}

type Handler struct {
	*Services
}

func New(services *Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func (h *Handler) Status(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func (h *Handler) JSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func (h *Handler) ParseBody(r *http.Request, v any) error {
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(v)

	var synerr *json.SyntaxError

	if err != nil && errors.As(err, &synerr) {
		return fmt.Errorf("malformed json at %d", synerr.Offset)
	}

	var typeerr *json.UnmarshalTypeError

	if err != nil && errors.As(err, &typeerr) {
		if funclib.NonWhiteSpace(typeerr.Field) {
			return fmt.Errorf("malformed json at key %s", typeerr.Field)
		}

		return fmt.Errorf("malformed json at %d", typeerr.Offset)
	}

	if err != nil && errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("malformed json")
	}

	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Errorf("request body has no content")
	}

	if err != nil {
		return err
	}

	return nil
}

type validator interface {
	Validate() error
}

func BIND[T validator](h *Handler, r *http.Request) (T, error) {
	var t T

	if err := h.ParseBody(r, &t); err != nil {
		return *new(T), err
	}

	if err := t.Validate(); err != nil {
		return *new(T), err
	}

	return t, nil
}

type handlerfn func(http.ResponseWriter, *http.Request) error

func wrap(fn handlerfn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)

		var he *HandlerError

		if err != nil && errors.As(err, &he) {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(he.code)
			w.Write([]byte(fmt.Sprintf(`{"detail": %s}`, strconv.Quote(he.msg))))
			return
		}

		if err != nil {
			panic(err)
		}
	}
}
