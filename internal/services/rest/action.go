package rest

import (
	"net/http"
)

func (s *service) Action(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
