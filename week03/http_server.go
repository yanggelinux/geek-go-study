package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var exit = make(chan struct{})

func helloWord(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello word!\n")
}

func shutDown(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "shut down!\n")
	exit <- struct{}{}
}

func main() {

	g, ctx := errgroup.WithContext(context.Background())
	//g := new(errgroup.Group)

	mux := http.NewServeMux()

	mux.HandleFunc("/hello", helloWord)
	mux.HandleFunc("/shutdown", shutDown)

	server := http.Server{
		Handler: mux,
		Addr:    ":8090",
	}

	g.Go(func() error {
		err := server.ListenAndServe()
		return err
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit...")
		case <-exit:
			log.Println("server eixt...")
		}

		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("shutting down server...")
		return server.Shutdown(timeoutCtx)
	})

	g.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.Errorf("get os signal: %v", sig)
		}
	})
	log.Println("server start...")
	err := g.Wait()
	if err != nil {
		log.Fatal("http server exit...:", err)
	}

}
