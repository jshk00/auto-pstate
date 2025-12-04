package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type Server struct {
	autoepp *AutoEPPSetter
	perfs   map[string]bool
	mux     *http.ServeMux
	srv     *http.Server
}

func NewServer(addr string, autoepp *AutoEPPSetter, perfs []string) *Server {
	m := make(map[string]bool)
	for _, p := range perfs {
		m[strings.TrimRight(p, "\n")] = true
	}
	mux := http.NewServeMux()
	return &Server{
		autoepp: autoepp,
		perfs:   m,
		mux:     mux,
		srv: &http.Server{
			Handler: mux,
			Addr:    addr,
		},
	}
}

func (s *Server) handleSetEPP(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("[INFO] profile %s set in manual mode\n", profile)
}

func (s *Server) handleAutoON(w http.ResponseWriter, r *http.Request) {
	m := r.PathValue("mode")
	mode := manual
	if m != "on" && m != "off" {
		http.Error(w, "invalid option auto path param should be on/off", http.StatusBadRequest)
		return
	}
	if m == "on" {
		mode = auto
		if err := FirstBoot(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	s.autoepp.SetMode(mode)
	log.Printf("[INFO] mode set to %s\n", mode)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	mode := s.autoepp.GetMode().String()
	profile, err := GetEPP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] status mode:%s, profile:%s\n", mode, profile)
	json.NewEncoder(w).Encode(map[string]any{
		"mode":    mode,
		"profile": profile,
	})
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /set/{profile}", s.handleSetEPP)
	s.mux.HandleFunc("POST /auto/{mode}", s.handleAutoON)
	s.mux.HandleFunc("POST /status", s.handleStatus)
}

func (s *Server) Start() error {
	s.routes()
	ln, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: SockPath,
		Net:  "unix",
	})
	if err != nil {
		return fmt.Errorf("failed to start unix listener: %w", err)
	}
	if err := s.srv.Serve(ln); err != nil {
		log.Println(err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
