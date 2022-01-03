// Package srvgroup implements a server runner with deterministic teardown.
package srvgroup

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultTimeout = 10 * time.Second

var defaultSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

// Group represents a server group
type Group struct {
	Signals []os.Signal
	Timeout time.Duration
}

// Run runs servers with a group with default signals and timeout
func Run(servers ...Server) []error {
	g := Group{
		Signals: defaultSignals,
		Timeout: defaultTimeout,
	}
	return g.Run(servers...)
}

// Run runs all servers concurrently.
// When the first Server.serve returns, all others will
// be shutdown. Returns only when all servers have finished serving.
// Returns all errors occurred while serving/shutting down
func (g Group) Run(servers ...Server) []error {
	var errors []error

	// nothing to do
	if len(servers) == 0 {
		return []error{}
	}

	serverDone := make(chan error, len(servers))

	for _, server := range servers {
		s := server
		go func() {
			serverDone <- s.Serve()
		}()
	}

	interrupted := make(chan os.Signal, 1)
	signal.Notify(interrupted, g.Signals...)

	// wait for the first server to stop or
	// for the user interruption
	select {
	case err := <-serverDone:
		errors = append(errors, err)
	case <-interrupted:
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()

	for _, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			errors = append(errors, err)
		}
	}

	// wait for all servers to stop (-1 because the first server is already stopped)
	for i := 1; i < cap(serverDone); i++ {
		err := <-serverDone
		if err != nil {
			// prepend the errors since they happened before shutdown
			errors = append([]error{err}, errors...)
		}
	}

	return errors
}
