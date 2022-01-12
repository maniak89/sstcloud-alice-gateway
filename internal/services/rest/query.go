package rest

import (
	"net/http"
)

func (s *service) Query(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
