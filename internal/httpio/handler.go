package httpio

import (
	"net/http"
	"short/linkit"
)

type Handler func(w http.ResponseWriter, r *http.Request) http.Handler

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if next := h(w, r); next != nil {
		next.ServeHTTP(w, r)
	}
}

func Error(code int, message string) Handler {
	return func(w http.ResponseWriter, r *http.Request) http.Handler {
		if code == http.StatusInternalServerError {
			Log(r.Context(), "%s: %v", r.URL.Path, message)
			message = linkit.ErrInternal.Error()
		}
		return JSON(code, map[string]string{
			"error": message,
		})
	}
}

func JSON(code int, v any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := Encode(w, code, v); err != nil {
			Log(r.Context(), "%s: JSON.Encode: %v", r.URL.Path, err)
		}
	}
}
