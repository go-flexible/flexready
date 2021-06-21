package flexready

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path"
)

const (
	// DefaultAddr is the port that we listen to the prometheus path on by default.
	DefaultAddr = "0.0.0.0:3674"

	// DefaultPath is the path where we expose prometheus by default.
	DefaultPath = "/ready"
)

// Config represents the configuration for the metrics server.
type Config struct {
	Path   string
	Server *http.Server
}

// New creates a new default metrics server.
func New(config *Config, checks Checks) *Server {
	if config == nil {
		config = &Config{}
	}
	if readinessPath := os.Getenv("READINESS_PATH"); readinessPath != "" && config.Path == "" {
		config.Path = readinessPath
	}
	if config.Path == "" {
		config.Path = DefaultPath
	}
	if config.Server == nil {
		config.Server = &http.Server{}
	}
	if addr := os.Getenv("READINESS_ADDR"); addr != "" && config.Server.Addr == "" {
		config.Server.Addr = addr
	}
	if config.Server.Addr == "" {
		config.Server.Addr = DefaultAddr
	}
	config.Server.Handler = CheckHandler(checks)
	return &Server{
		Checks: checks,
		Server: config.Server,
		Path:   path.Join("/", config.Path),
	}
}

// Server defines a readiness server.
type Server struct {
	*http.Server

	Path   string
	Checks Checks
}

// Run will start the ready server.
func (s *Server) Run(_ context.Context) error {
	lis, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return err
	}
	log.Printf("serving readiness checks server over http on http://%s%s", s.Addr, s.Path)
	return s.Server.Serve(lis)
}

// Halt will attempt to gracefully shut down the server.
func (s *Server) Halt(ctx context.Context) error {
	log.Printf("stopping readiness checks server over http on http://%s...", s.Addr)
	return s.Server.Shutdown(ctx)
}

// CheckHandler provides a function for providing health checks over http.
func CheckHandler(checks Checks) http.HandlerFunc {
	type health struct {
		OK      bool   `json:"ok"`
		Message string `json:"messages"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var ready = true
		res := make(map[string]health)
		for name, check := range checks {
			var message string
			err := check.Check()
			if err != nil {
				ready = false
				message = err.Error()
			}
			res[name] = health{
				OK:      err == nil,
				Message: message,
			}
		}

		w.Header().Add("Content-Type", "application/json")
		bts, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		code := http.StatusOK
		if !ready {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		_, _ = w.Write(bts)
	}
}
