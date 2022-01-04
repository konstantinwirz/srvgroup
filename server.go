package srvgroup

import (
	"context"
	"net/http"
)

type (

	// Server represents a unit which lifecycle
	// is managed by the Group
	Server struct {
		Serve    func() error
		Shutdown func(ctx context.Context) error
	}

	// ServerMiddleware is a middleware for the Server
	ServerMiddleware = func(Server) Server

	ServerLifecycleHooks struct {
		BeforeServe    func()
		AfterServe     func(error)
		BeforeShutdown func()
		AfterShutdown  func(error)
	}
)

// ServerLifecycleMiddleware creates a middleware that allows
// to create lifecycle hooks which will be executed before and after Serve
// and Shutdown functions.
func ServerLifecycleMiddleware(hooks ServerLifecycleHooks) ServerMiddleware {
	return func(next Server) Server {
		return Server{
			Serve: func() error {
				if hooks.BeforeServe != nil {
					hooks.BeforeServe()
				}

				err := next.Serve()

				if hooks.AfterServe != nil {
					hooks.AfterServe(err)
				}

				return err
			},
			Shutdown: func(ctx context.Context) error {
				if hooks.BeforeShutdown != nil {
					hooks.BeforeShutdown()
				}

				err := next.Shutdown(ctx)

				if hooks.AfterShutdown != nil {
					hooks.AfterShutdown(err)
				}

				return err
			},
		}
	}
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
