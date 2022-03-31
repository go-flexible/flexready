package flexready

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultAddr is the port that we expose the readiness server on by default.
	DefaultAddr = "0.0.0.0:3674"

	// DefaultReadinessPath is the path where we readiness by default.
	DefaultReadinessPath = "/ready"

	// DefaultLivenessPath is the path where we expose liveness by default.
	DefaultLivenessPath = "/live"

	// DefaultReadTimeout is the default read timeout for the http server.
	DefaultReadTimeout = 5 * time.Second

	// DefaultReadHeaderTimeout is the default read header timeout for the http server.
	DefaultReadHeaderTimeout = 1 * time.Second

	// DefaultIdleTimeout is the default idle timeout for the http server.
	DefaultIdleTimeout = 1 * time.Second

	// DefaultWriteTimeout is the default write timeout for the http server.
	DefaultWriteTimeout = 15 * time.Second
)

var defaultLogger = log.New(os.Stderr, "flexready: ", 0)

// Option is a type of func that allows you change defaults of the *Server
// returned by New.
type Option func(s *Server)

// WithLogger allows you to set a logger for the server.
func WithLogger(l Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}

// WithLivenessPath allows you to set the path for the liveness endpoint.
func WithLivenessPath(path string) Option {
	return func(s *Server) {
		s.livenessPath = path
	}
}

// WithReadinessPath allows you to set the path for the readiness endpoint.
func WithReadinessPath(path string) Option {
	return func(s *Server) {
		s.readinessPath = path
	}
}

// WithAddress allows you to set the address for the server.
func WithAddress(address string) Option {
	return func(s *Server) {
		s.address = address
	}
}

// WithHTTPServer allows you to set the http server for the server.
func WithHTTPServer(s *http.Server) Option {
	return func(ss *Server) {
		// prefer the address of the provided server.
		ss.address = s.Addr
		ss.Server = s
	}
}

// New creates a new ready server.
func New(checks Checks, options ...Option) *Server {
	var (
		address       = DefaultAddr
		livenessPath  = DefaultLivenessPath
		readinessPath = DefaultReadinessPath
	)
	// if defined use the env vars.
	if ad := os.Getenv("FLEX_READYSRV_ADDR"); ad != "" {
		address = ad
	}
	if lp := os.Getenv("FLEX_READYSRV_LIVENESS_PATH"); lp != "" {
		livenessPath = lp
	}
	if rp := os.Getenv("FLEX_READYSRV_READINESS_PATH"); rp != "" {
		readinessPath = rp
	}

	server := &Server{
		logger:        defaultLogger,
		address:       address,
		livenessPath:  livenessPath,
		readinessPath: readinessPath,
		checks:        checks,
		Server: &http.Server{
			ReadTimeout:       DefaultReadTimeout,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			IdleTimeout:       DefaultIdleTimeout,
		},
	}

	for _, option := range options {
		option(server)
	}

	mux := http.NewServeMux()
	mux.Handle(server.livenessPath, LivenessHandler())
	mux.Handle(server.readinessPath, ReadinessHandler(checks))

	server.Server.Addr = server.address
	server.Server.Handler = mux

	return server
}

// Logger defines any logger able to call Printf.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Server defines a readiness server.
type Server struct {
	*http.Server
	logger        Logger
	checks        Checks
	address       string
	livenessPath  string
	readinessPath string
}

// Run will start the ready server.
func (s *Server) Run(_ context.Context) error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	if address, ok := listener.Addr().(*net.TCPAddr); ok {
		s.logger.Printf("serving readiness checks over http on http://%s%s", address, s.readinessPath)
		s.logger.Printf("serving liveness checks over http on http://%s%s", address, s.livenessPath)
	}
	return s.Serve(listener)
}

// Halt will attempt to gracefully shut down the server.
func (s *Server) Halt(ctx context.Context) error {
	s.logger.Printf("shutting down readiness server...")
	return s.Server.Shutdown(ctx)
}

// LivenessHandler is a function for providing liveness checks over http.
func LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// ReadinessHandler is a function for providing health checks over http.
func ReadinessHandler(checks Checks) http.HandlerFunc {
	type health struct {
		Message string `json:"messages"`
		OK      bool   `json:"ok"`
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		ready := true
		res := make(map[string]health)
		for name, check := range checks {
			var message string
			err := check()
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
