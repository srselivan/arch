package http

import (
	"arch/internal/server/controller/http/middleware"
	v1 "arch/internal/server/controller/http/v1"
	"arch/internal/server/entity"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"net/http"
	"os"
)

type messageService interface {
	Save(ctx context.Context, message entity.Message)
	Get(cxt context.Context, id string) (entity.Message, error)
	GetAll(ctx context.Context) ([]entity.Message, error)
}

type authService interface {
	Login(ctx context.Context, credentials entity.UserCredentials, resource entity.ResourceInfo) error
}

type Config struct {
	Addr string
	Log  *zerolog.Logger

	MessageService messageService
	AuthService    authService
}

type Server struct {
	addr       string
	httpServer *http.Server
	log        *zerolog.Logger

	messageService messageService
	authService    authService
}

func New(cfg *Config) *Server {
	s := &Server{
		addr:           cfg.Addr,
		httpServer:     nil,
		log:            cfg.Log,
		messageService: cfg.MessageService,
		authService:    cfg.AuthService,
	}

	s.httpServer = &http.Server{
		Addr: s.addr,
	}
	s.init()

	return s
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) init() {
	v1Handler := v1.New(v1.Config{
		Log:            s.log,
		MessageService: s.messageService,
	})

	r := chi.NewRouter()

	dir, _ := os.Getwd()
	r.Handle(
		"/_upload/*",
		http.StripPrefix(
			"/_upload",
			http.FileServer(http.Dir(dir+"/_upload")),
		),
	)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.NewBasicAuth(s.authService, s.log))
		r.Mount("/messages", v1Handler.InitMessageRoutes())
		r.Mount("/files", v1Handler.InitFileRoutes())
	})
	r.Mount("/", v1Handler.InitWebRoutes())

	s.httpServer.Handler = r
}
