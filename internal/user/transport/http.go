package transport

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ribice/chisk/internal/user"
)

// New instantates user http transport
func New(r *chi.Mux, svc *user.Service) {
}

// Service represents user http service
type Service struct {
	svc *user.Service
}

func (s *Service) create(w http.ResponseWriter, r *http.Request) {
}
