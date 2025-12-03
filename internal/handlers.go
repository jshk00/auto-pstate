package internal

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	autoepp *AutoEPPSetter
	perfs   map[string]bool
	mux     *http.ServeMux
}

func NewServer(autoepp *AutoEPPSetter, perfs []string) *Server {
	srv := &Server{autoepp: autoepp, perfs: make(map[string]bool), mux: http.NewServeMux()}
	for _, p := range perfs {
		srv.perfs[p] = true
	}
	srv.mux.HandleFunc("POST /set/{profile}", srv.HandleSetEPP)
	srv.mux.HandleFunc("POST /auto/on", srv.HandleAutoON)
	srv.mux.HandleFunc("POST /status", srv.HandleStatus)
	return srv
}

func (s *Server) HandleSetEPP(w http.ResponseWriter, r *http.Request) {
	profile := r.PathValue("profile")
	if !s.perfs[profile] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.autoepp.SetMode(manual)
	if err := SetEPP(profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleAutoON(w http.ResponseWriter, _ *http.Request) {
	s.autoepp.SetMode(auto)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleStatus(w http.ResponseWriter, _ *http.Request) {
	mode := s.autoepp.GetMode().String()
	profile, err := GetEPP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"current_mode":    mode,
		"current_profile": profile,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
