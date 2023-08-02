package short

import (
	"fmt"
	"net/http"
	"short/internal/httpio"
	"short/linkit"
)

const (
	shorteningRoute  = "/s"
	resolveRoute     = "/r/"
	healthCheckRoute = "/health"
)

type mux http.Handler

type Server struct {
	mux
}

func NewServer() *Server {
	var s Server
	s.registerRoutes()
	return &s
}

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc(shorteningRoute, s.shorteningHandler)
	mux.HandleFunc(resolveRoute, s.resolveHandler)
	mux.HandleFunc(healthCheckRoute, s.healthCheckHandler)
	s.mux = mux
}

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func (s *Server) shorteningHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		URL string
		Key string
	}
	err := httpio.Decode(http.MaxBytesReader(w, r.Body, 4_096), &input)
	if err != nil {
		http.Error(w, "cannot decode JSON", http.StatusBadRequest)
		return
	}
	ln := link{
		uri:      input.URL,
		shortKey: input.Key,
	}
	if err := checkLink(ln); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = httpio.Encode(w, http.StatusCreated, map[string]any{
		"key": ln.shortKey,
	})
}

func (s *Server) resolveHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(resolveRoute):]

	if err := checkShortKey(key); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// use dummy data for now and carelessly expose internal details.
	if key == "fortesting" {
		http.Error(w, "db at IP ... failed", http.StatusInternalServerError)
		return
	}
	if key != "go" {
		http.Error(w, linkit.ErrNotExist.Error(), http.StatusNotFound)
		return
	}
	const uri = "https://go.dev"
	http.Redirect(w, r, uri, http.StatusFound)
}
