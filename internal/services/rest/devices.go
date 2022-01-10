package rest

import (
	"net/http"
)

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
