package handler

import "net/http"

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) error {
	return h.JSON(w, 200, map[string]string{"status": "ok"})
}
