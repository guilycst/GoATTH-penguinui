package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/guilycst/gottha-penguinui/internal/pages/demo/components"
)

// Server handles HTTP requests for both original PenguinUI and GoTTHA components
type Server struct {
	projectRoot string
	mux         *http.ServeMux
}

// New creates a new server instance
func New(projectRoot string) *Server {
	s := &Server{
		projectRoot: projectRoot,
		mux:         http.NewServeMux(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	// Compiled assets (CSS, JS)
	assetsDir := filepath.Join(s.projectRoot, "assets")
	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))
	s.mux.Handle("/assets/", assetsHandler)

	// Original PenguinUI static files (for direct access if needed)
	originalDir := filepath.Join(s.projectRoot)
	originalHandler := http.StripPrefix("/original/", http.FileServer(http.Dir(originalDir)))
	s.mux.Handle("/original/", originalHandler)

	// Component comparison pages
	s.mux.HandleFunc("/components/", s.handleComponent)

	// API endpoints for HTMX demos
	s.mux.HandleFunc("/api/hello", s.handleAPIHello)
	s.mux.HandleFunc("/api/components/button", s.handleButtonFragment)

	// Root redirect to first component
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/components/button", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})
}

func (s *Server) handleComponent(w http.ResponseWriter, r *http.Request) {
	// Extract component name from URL
	path := strings.TrimPrefix(r.URL.Path, "/components/")
	if path == "" {
		http.Redirect(w, r, "/components/button", http.StatusMovedPermanently)
		return
	}

	componentName := strings.Split(path, "/")[0]

	switch componentName {
	case "button":
		s.renderButtonComponent(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) renderButtonComponent(w http.ResponseWriter, r *http.Request) {
	// Load original HTML content
	originalHTML, err := s.loadOriginalHTML("buttons/default-button.html")
	if err != nil {
		// If file doesn't exist, use a placeholder
		originalHTML = "<!-- Original HTML not found -->"
	}

	// Render the component demo page
	components.ButtonDemoPage(originalHTML).Render(r.Context(), w)
}

func (s *Server) loadOriginalHTML(filename string) (string, error) {
	filepath := filepath.Join(s.projectRoot, filename)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *Server) handleAPIHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<p class="text-green-600">Hello from HTMX! Request received at %s %s</p>`, r.Method, r.URL.Path)
}

func (s *Server) handleButtonFragment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Get disabled state from query params
	disabled := r.URL.Query().Get("disabled") == "true"

	// Render just the button grid fragment
	components.ButtonFragment(disabled).Render(r.Context(), w)
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
