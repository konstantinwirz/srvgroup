package srvgroup

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestServerLifecycleMiddleware(t *testing.T) {
	expectedServeErr := errors.New("serve error")
	expectedShutdownErr := errors.New("shutdown error")

	beforeServe := make(chan struct{}, 1)
	afterServe := make(chan error, 1)
	beforeShutdown := make(chan struct{}, 1)
	afterShutdown := make(chan error, 1)
	done := make(chan struct{}, 1)

	var mu sync.Mutex

	srv := ServerLifecycleMiddleware(
		ServerLifecycleHooks{
			BeforeServe: func() {
				mu.Lock()
				beforeServe <- struct{}{}
			},
			AfterServe: func(err error) {
				mu.Lock()
				afterServe <- err
			},
			BeforeShutdown: func() {
				mu.Lock()
				beforeShutdown <- struct{}{}
			},
			AfterShutdown: func(err error) {
				mu.Lock()
				afterShutdown <- err
			},
		},
	)(
		Server{
			Serve: func() error {
				return expectedServeErr
			},
			Shutdown: func(ctx context.Context) error {
				return expectedShutdownErr
			},
		},
	)

	go func() {
		Run(srv)
		done <- struct{}{}
	}()

	// check if beforeServe was called
	select {
	case <-beforeServe:
		mu.Unlock()
	case <-time.After(10 * time.Millisecond):
		t.Fatal("beforeServe was not called")
	}

	// check if afterServe was performed
	select {
	case err := <-afterServe:
		if !reflect.DeepEqual(err, expectedServeErr) {
			t.Fatalf("expected %v, got %v", expectedServeErr, err)
		}
		mu.Unlock()
	case <-time.After(10 * time.Millisecond):
		t.Fatal("afterServe ws not called")
	}

	// check if beforeShutdown was performed
	select {
	case <-beforeShutdown:
		mu.Unlock()
	case <-time.After(10 * time.Millisecond):
		t.Fatal("beforeShutdown was not called")
	}

	// check if afterShutdown was performed
	select {
	case err := <-afterShutdown:
		if !reflect.DeepEqual(err, expectedShutdownErr) {
			t.Fatalf("expected %v, got %v", expectedServeErr, err)
		}
		mu.Unlock()
	case <-time.After(10 * time.Millisecond):
		t.Fatal("afterShutdown was not called")
	}

	// check if Run function is done
	select {
	case <-done:
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Run was not completed")
	}
}
