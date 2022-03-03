package server

import (
	"context"
	"net/http"
	"time"

	"github.com/AliRostami1/baagh/pkg/logy"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer   *http.Server
	router       *mux.Router
	logger       logy.Logger
	ctx          context.Context
	shutdownWait time.Duration
}

func New(ctx context.Context, log logy.Logger, addr string, writeTimeout time.Duration, readTimeout time.Duration, idleTimeout time.Duration, shutdownWait time.Duration) *Server {
	r := mux.NewRouter()
	// Add your routes as needed

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  time.Second * readTimeout,
			WriteTimeout: time.Second * writeTimeout,
			IdleTimeout:  time.Second * idleTimeout,
		},
		router:       r,
		logger:       log,
		ctx:          ctx,
		shutdownWait: shutdownWait,
	}
}

func (s *Server) initMiddlewares() {
	s.router.Use(s.loggingMiddleware)
}

func (s *Server) initRoutes() {
	apiSub := s.router.PathPrefix("/api").Subrouter()

	chipsSub := apiSub.PathPrefix("/chips").Subrouter()

	chipsSub.HandleFunc("/", s.getAllChips).Methods("GET")
	chipsSub.HandleFunc("/{chip}", s.getOneChip).Methods("GET")

	itemsSub := chipsSub.PathPrefix("/{chip}/items").Subrouter()

	itemsSub.HandleFunc("/", s.getAllItems).Methods("GET")
	// itemsSub.HandleFunc("/", itemsWatchHandler).Methods("GET").Queries("watch", "true")
	// itemsSub.HandleFunc("/", itemsPostHandler).Methods("POST")
	// itemsSub.HandleFunc("/", itemsDeleteHandler).Methods("DELETE")

	itemsSub.HandleFunc("/{offset:[0-9]+}", s.getOneItem).Methods("GET")
	// itemsSub.HandleFunc("/{offset:[0-9]+}", itemWatchHandler).Methods("GET").Queries("watch", "true")
	itemsSub.HandleFunc("/{offset:[0-9]+}", s.createOneItem).Methods("POST")
	itemsSub.HandleFunc("/{offset:[0-9]+}", s.deleteOneItem).Methods("DELETE")

	// apiSub.HandleFunc("/healthcheck", healthCheckHandler).Methods("GET")

	apiSub.HandleFunc("/version", versionHandler).Methods("GET")
}

func (s *Server) Start() error {
	// initialize middlewares
	s.initMiddlewares()

	// initialize routes
	s.initRoutes()

	// start the server in another goroutine
	go func() {
		s.logger.Infof("starting the server on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
			s.logger.Infof("server stopped: %v", err)
		} else {
			s.logger.Errorf("server crashed: %v", err)
		}
	}()

	// wait for ctx to close
	<-s.ctx.Done()

	// gracefully shutdown the server, waiting at most
	// for shutdownWait duration
	toCtx, cancel := context.WithTimeout(context.Background(), s.shutdownWait)
	defer cancel()

	return s.httpServer.Shutdown(toCtx)
}
