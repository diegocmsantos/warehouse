package whhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	warehouse "github.com/diegocmsantos/warehouse/core"
	"github.com/diegocmsantos/warehouse/db"
	"go.uber.org/zap"
)

const shutdownTimeout = time.Minute

// Server represents a HTTP server
type Server struct {
	httpServer        *http.Server
	mux               *http.ServeMux
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	host              string
	port              int
	logger            *zap.Logger
}

// ServerOptions represents the options for creating a server
type ServerOptions struct {
	Addr        string
	Host        string
	Port        int
	ReadTimeout time.Duration
	Logger      *zap.Logger
}

// New creates a new server
func New(opts ServerOptions) (*Server, error) {
	return &Server{
		host:   opts.Host,
		port:   opts.Port,
		mux:    http.NewServeMux(),
		logger: opts.Logger,
	}, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.HandleIndex())
}

// Start starts the Server.
func (s *Server) Start() error {
	s.routes()

	address := net.JoinHostPort(s.host, strconv.Itoa(s.port))
	s.httpServer = &http.Server{
		Addr:              address,
		Handler:           s.mux,
		ReadTimeout:       s.readTimeout,
		ReadHeaderTimeout: s.readHeaderTimeout,
		WriteTimeout:      s.writeTimeout,
		IdleTimeout:       s.idleTimeout,
	}

	s.logger.Sugar().Infof("Starting on %v", address)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error serving http: %+v", err)
	}

	return nil
}

// Stop stops the Server.
func (s *Server) Stop() error {
	s.logger.Info("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isCSVContent := r.Header.Get("Content-Type") == "text/csv"
		if r.Method == http.MethodPost && isCSVContent {

			defer r.Body.Close()

			csvSource, err := db.NewWithReader(r.Body, db.CategoryHiearchy)
			if err != nil {
				s.logger.Error("error creating csv source service: ", zap.Error(err))
				Error(w, err)
				return
			}

			wh, err := warehouse.New(csvSource)
			if err != nil {
				s.logger.Error("error creating warehouse service: ", zap.Error(err))
				Error(w, err)
				return
			}

			categoryList, err := wh.CreateList()
			if err != nil {
				s.logger.Error("error reading warehouse service: ", zap.Error(err))
				code := http.StatusInternalServerError
				var hierarchyErr warehouse.ErrCategoryHiearchy
				if errors.As(err, &hierarchyErr) {
					code = http.StatusBadRequest
				}
				ErrorWithCode(w, code, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(categoryList)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
