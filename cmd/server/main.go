package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jshk00/auto-pstate/internal"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := internal.IsRoot(); err != nil {
		return err
	}
	if err := internal.IsPState(); err != nil {
		return err
	}
	if err := internal.SetGovernor(); err != nil {
		return err
	}
	if err := internal.FirstBoot(); err != nil {
		return err
	}

	perfs, err := internal.GetPreferences()
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	autoepp := &internal.AutoEPPSetter{}
	srv := internal.NewServer(":9010", autoepp, perfs)

	// Closer routine
	go func() {
		<-ch
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		autoepp.Close()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println(err)
		}
	}()

	autoepp.Start()
	return srv.Start()
}
