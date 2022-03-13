package server

import (
	"context"

	"github.com/Kran001/basic-auth/internal/store"

	"github.com/Kran001/basic-auth/pkg/httpsserver"
	"github.com/Kran001/basic-auth/pkg/logging"

	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	*httpsserver.Server
	repos *store.Repositories

	certPath string
	keyPath  string
}

func NewServer(repos *store.Repositories,
	config httpsserver.Config,
	certPath,
	keyPath string,
) *Server {
	return &Server{
		Server:   httpsserver.NewServer(config),
		repos:    repos,
		certPath: certPath,
		keyPath:  keyPath,
	}
}

func (s *Server) Init() error {
	c := cors.AllowAll()

	s.SetHTTPSRouter(c.Handler(s.initRoutes()))

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	logging.Logger.Info("Generate cert for https server")
	err := httpsserver.CheckCerts(s.certPath, s.keyPath, s.HTTPSConnectionString())
	if err != nil {
		logging.Logger.Error("Error: can not generate certificate.")
	}

	logging.Logger.Info("Start API server handlers")

	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return s.StartHTTPSServe(s.certPath, s.keyPath)
	})

	logging.Logger.Info("The https and http services are ready to listen and serve.")

	return eg.Wait()
}
