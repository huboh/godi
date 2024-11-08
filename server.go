package godi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	mux    *http.ServeMux
	server *http.Server
}

func newHttpServer(mux *http.ServeMux) *HttpServer {
	return &HttpServer{
		mux: mux,
		server: &http.Server{
			Handler: mux,
		},
	}
}

// Listen starts the HTTP server on the specified host and port, and listens
// for incoming requests.
//
// Also listens for system signals like SIGINT and SIGTERM to enable graceful shutdown.
func (s *HttpServer) Listen(host string, port string) error {
	errChan := make(chan error, 1)
	sigChan := make(chan os.Signal, 1)
	s.server.Addr = net.JoinHostPort(host, port)

	// listen for signals to allow graceful shutdown
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer close(errChan)

		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	log.Printf("listening on (%s)\n", s.server.Addr)

	select {
	case <-sigChan:
		return s.Shutdown(context.Background())

	case err := <-errChan:
		return fmt.Errorf("error listening on (%s) : %w", s.server.Addr, err)
	}
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HttpServer) Shutdown(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, (time.Second * 5))
	defer cancel()
	return s.server.Shutdown(ctx)
}
