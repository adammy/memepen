package meme

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adammy/memepen-services/pkg/httpapi"
	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	config *Config
	router *chi.Mux
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: chi.NewRouter(),
	}
}

func (s *Server) Start() error {
	openapi, err := GetSwagger()
	if err != nil {
		return err
	}
	openapi.Servers = nil

	logger := httpapi.NewLogger(s.config.Logger)

	s.router.Use(oapimiddleware.OapiRequestValidatorWithOptions(openapi, &oapimiddleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			httpapi.SendErrorJSON(w, statusCode, errors.New(message))
		},
	}))
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(s.config.RequestTimeout))
	s.router.Use(httplog.RequestLogger(logger))

	HandlerFromMux(s, s.router)

	logger.Info().Msgf("server starting on port %d", s.config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router); err != nil {
		return err
	}

	return nil
}

func (s *Server) CreateMeme(w http.ResponseWriter, r *http.Request) {
	// implement me
}

func (s *Server) GetMeme(w http.ResponseWriter, r *http.Request, memeID string) {
	// implement me
}