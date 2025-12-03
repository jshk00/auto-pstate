package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jshk00/auto-pstate/internal"
)

func main() {
	log.SetFlags(0)
	if err := internal.IsRoot(); err != nil {
		log.Fatal(err)
	}
	if err := internal.IsPState(); err != nil {
		log.Fatal(err)
	}
	if err := internal.SetGovernor(); err != nil {
		log.Fatal(err)
	}

	perfs, err := internal.GetPreferences()
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("/run/pstated/pstated.sock", ":9010")
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	autoepp := &internal.AutoEPPSetter{}
	go autoepp.Run()
	defer autoepp.Close()

	srv := &http.Server{
		Handler: internal.NewServer(autoepp, perfs),
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			log.Fatal(err)
		}
	}()
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
