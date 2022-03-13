package server

import (
	"github.com/Kran001/basic-auth/pkg/httpsserver"

	"github.com/gorilla/mux"
)

func (s *Server) initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(httpsserver.Recover)
	r.Use(httpsserver.LogRequest)

	s.initUsersRoutes(r)
	s.initSessionsRoutes(r)

	return r
}
