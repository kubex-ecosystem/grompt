package server

import (
    "net/http"
)

// withAPI wraps a handler with API concerns: CORS headers, OPTIONS preflight,
// and optional method allowlisting. It delegates CORS header values to
// Handlers.setCORSHeaders to keep a single source of truth.
func withAPI(h *Handlers, fn http.HandlerFunc, methods ...string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if h != nil {
            h.setCORSHeaders(w)
        }
        if r.Method == http.MethodOptions {
            return
        }
        if len(methods) > 0 {
            ok := false
            for _, m := range methods {
                if r.Method == m {
                    ok = true
                    break
                }
            }
            if !ok {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
            }
        }
        fn(w, r)
    }
}

// routeBuilder provides a tiny fluent API for registering routes with options.
type routeBuilder struct {
    s       *Server
    path    string
    handler http.HandlerFunc
    methods []string
    useAPI  bool
}

// Route creates a new routeBuilder for a given path and handler.
func (s *Server) Route(path string, h http.HandlerFunc) *routeBuilder {
    return &routeBuilder{s: s, path: path, handler: h}
}

// Methods sets the allowed HTTP methods for the route.
func (rb *routeBuilder) Methods(m ...string) *routeBuilder {
    rb.methods = append([]string{}, m...)
    return rb
}

// WithAPI enables API wrapper (CORS/OPTIONS/methods) for the route.
func (rb *routeBuilder) WithAPI() *routeBuilder {
    rb.useAPI = true
    return rb
}

// Register registers the route in the underlying mux with all configured options.
func (rb *routeBuilder) Register() {
    h := rb.handler
    if rb.useAPI {
        h = withAPI(rb.s.handlers, h, rb.methods...)
    }
    rb.s.router.HandleFunc(rb.path, h)
}

