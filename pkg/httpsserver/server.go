package httpsserver

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	log "github.com/Kran001/basic-auth/pkg/logging"
)

type Server struct {
	Config

	serverHTTPS *http.Server
}

func NewServer(config Config) *Server {
	s := &Server{
		Config: config,
		serverHTTPS: &http.Server{
			Addr: config.HTTPSConnectionString(),
		},
	}

	return s
}

func (s *Server) SetHTTPSRouter(handler http.Handler) {
	s.serverHTTPS.Handler = handler
}

func (s *Server) StartHTTPSServe(cert, key string) error {
	log.Logger.Info("HTTPS Server starting")
	return s.serverHTTPS.ListenAndServeTLS(cert, key)
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Logger.Info("Shutdown server")
	if err := s.serverHTTPS.Shutdown(ctx); err != nil {
		log.Logger.Errorf("Error while shutdown Server: %s\n", err.Error())
	}
}
