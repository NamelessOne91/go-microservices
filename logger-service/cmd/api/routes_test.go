package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func routeExists(routes chi.Router, route string) bool {
	found := false
	_ = chi.Walk(routes, func(method, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	})
	return found
}

func Test_routes_exist(t *testing.T) {
	tests := []struct {
		name     string
		route    string
		expected bool
	}{
		{"Route should exist", "/log", true},
		{"Route does not exist", "/test-route", false},
	}

	testApp := Config{}
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := routeExists(chiRoutes, tt.route)
			if got != tt.expected {
				t.Errorf("expected routeExists('%s') to be %t - got %t", tt.route, tt.expected, got)
			}
		})
	}

}
