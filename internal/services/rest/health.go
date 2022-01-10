package rest

import (
	"net/http"
)

func (s *service) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
