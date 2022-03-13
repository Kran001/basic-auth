package server

import (
	"github.com/Kran001/basic-auth/pkg/logging"
	"net/http"
	"strconv"

	"github.com/Kran001/basic-auth/pkg/utils"

	"github.com/gorilla/mux"
)

func (s *Server) initSessionsRoutes(r *mux.Router) {
	r.HandleFunc("/sessions",
		utils.CheckToken(s.getAllUsersSessionsInfo)).Methods(http.MethodGet)
	r.HandleFunc("/sessions/{uid:[0-9]+}",
		utils.CheckToken(s.getUserSessions)).Methods(http.MethodGet)
	r.HandleFunc("/sessions/one/{sid:[0-9]+}",
		utils.CheckToken(s.deleteUserSessionById)).Methods(http.MethodDelete)
	r.HandleFunc("/sessions/all/{uid:[0-9]+}",
		utils.CheckToken(s.deleteAllUserSessions)).Methods(http.MethodDelete)
}

func (s *Server) getAllUsersSessionsInfo(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Info("getAllUsersSessionsInfo handler")
	logging.Logger.Info("Getting all users sessions info...")

	res, err := s.repos.Sessions.AllUsersSessionsInfo(request.Context())
	if err != nil {
		logging.Logger.Error("Failed load all users sessions list. Reason:", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Successfully loaded all users sessions list.")

	utils.Respond(writer, map[string]interface{}{"res": res}, http.StatusOK)
}

// Process GET req for getting all users sessions info
func (s *Server) getUserSessions(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("getUserSessions handler")
	logging.Logger.Info("Getting user sessions list...")

	vars := mux.Vars(request)
	userId, err := strconv.Atoi(vars["uid"])
	if err != nil {
		logging.Logger.Error("Failed parse the user ID, Reason:", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Trace("user ID is : ", userId)

	logging.Logger.Info("Load user sessions list")
	res, err := s.repos.Sessions.UserSessions(request.Context(), int64(userId))
	if err != nil {
		logging.Logger.Error("Failed load the user sessions list. Reason:", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Successfully loaded user sessions list.")

	utils.Respond(writer, map[string]interface{}{"res": res}, http.StatusOK)
}

// Process DELETE req for deleting ONE user session by user id and session id
func (s *Server) deleteUserSessionById(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("deleteUserSessionById handler")
	logging.Logger.Info("Deleting user session by id...")

	vars := mux.Vars(request)
	sessionId, err := strconv.Atoi(vars["sid"])
	if err != nil {
		logging.Logger.Error("Failed parse the user ID, Reason:", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Delete user session from store")
	err = s.repos.Sessions.DeleteUserSessionById(request.Context(), int64(sessionId))
	if err != nil {
		logging.Logger.Error("Failed deleting user session from store. Reason", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)
		return
	}

	logging.Logger.Info("Successfully deleted.")

	utils.Respond(writer, map[string]interface{}{"res": "ok"}, http.StatusOK)
}

// Process DELETE req for deleting ALL user session by user id
func (s *Server) deleteAllUserSessions(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("deleteAllUserSessions handler")
	logging.Logger.Info("Deleting all user sessions...")

	vars := mux.Vars(request)
	userId, err := strconv.Atoi(vars["uid"])
	if err != nil {
		logging.Logger.Error("Failed parse the user ID, Reason:", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Delete all user sessions from store")
	err = s.repos.Sessions.DeleteAllUserSessions(request.Context(), int64(userId))
	if err != nil {
		logging.Logger.Error("Failed deleting all user sessions from store. Reason", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Successfully deleted.")

	utils.Respond(writer, map[string]interface{}{"res": "ok"}, http.StatusOK)
}
