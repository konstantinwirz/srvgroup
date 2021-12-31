package srvgroup

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		servers    []Server
		wantErrors []error
	}{
		{
			name:       "empty servers",
			servers:    []Server{},
			wantErrors: []error{},
		},
		{
			name: "one server",
			servers: []Server{{
				Serve: func() error {
					return errors.New("serve error")
				},
				Shutdown: func(ctx context.Context) error {
					return errors.New("shutdown error")
				},
			}},
			wantErrors: []error{errors.New("serve error"), errors.New("shutdown error")},
		},
		{
			name: "many servers",
			servers: []Server{
				{
					Serve: func() error {
						return errors.New("serve error 1")
					},
					Shutdown: func(context.Context) error {
						return errors.New("shutdown error 1")
					},
				},
				{
					Serve: func() error {
						return errors.New("serve error 2")
					},
					Shutdown: func(context.Context) error {
						return errors.New("shutdown error 2")
					},
				},
				{
					Serve: func() error {
						return errors.New("serve error 3")
					},
					Shutdown: func(context.Context) error {
						return errors.New("shutdown error 3")
					},
				},
			},
			wantErrors: []error{errors.New("serve error 1"), errors.New("shutdown error 1"),
				errors.New("serve error 2"), errors.New("shutdown error 2"),
				errors.New("serve error 3"), errors.New("shutdown error 3")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorsChan := make(chan []error, 1)
			go func() {
				errorsChan <- Run(tt.servers...)
			}()

			select {
			case gotErrors := <-errorsChan:
				if len(gotErrors) != len(tt.wantErrors) {
					t.Fatalf("expected %d errors, got %d", len(tt.wantErrors), len(gotErrors))
				}

				for _, wantErr := range tt.wantErrors {
					if hasError(wantErr, gotErrors) {
						t.Fatalf("expected %v to be returned by Run", wantErr)
					}
				}

			case <-time.After(100 * time.Millisecond):
				t.Fatal("timeout")
			}
		})
	}
}

func hasError(err error, errors []error) bool {
	for _, e := range errors {
		if e == err {
			return true
		}
	}

	return false
}
