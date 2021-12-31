package srvgroup

import (
	"context"
	"net/http"
)

// Server represents a unit which lifecycle
// is managed by the Group
type Server struct {
	Serve    func() error
	Shutdown func(ctx context.Context) error
}

// HTTPServer makes a Server from given http.Server
func HTTPServer(srv *http.Server) Server {
	return Server{
		Serve: func() error {
			err := srv.ListenAndServe()
			if err == nil || err == http.ErrServerClosed {
				return nil
			}
			return err
		},
		Shutdown: func(ctx context.Context) error {
			if err := srv.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		},
	}
}
