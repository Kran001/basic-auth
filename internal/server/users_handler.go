package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Kran001/basic-auth/internal/domain"

	"github.com/Kran001/basic-auth/pkg/utils"
	"github.com/Kran001/basic-auth/pkg/validators"

	"github.com/Kran001/basic-auth/pkg/logging"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) initUsersRoutes(r *mux.Router) {
	r.HandleFunc("/register", s.registerRequest).Methods(http.MethodPost)
	r.HandleFunc("/login", s.loginRequest).Methods(http.MethodPost)
	r.HandleFunc("/logout", utils.CheckToken(s.logoutRequest)).Methods(http.MethodPost, http.MethodGet)
}

func (s *Server) registerRequest(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("registerRequest handler")
	logging.Logger.Info(request.RemoteAddr, " register...")

	user := domain.User{}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		logging.Logger.Error("Failed request decoding: ", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}
	logging.Logger.Info("Registering new user as local.")

	logging.Logger.Debug("Validate logging stage starting...")
	valid := validators.ValidateByPattern(user.Name, validators.LoginPattern)
	if !valid {
		logging.Logger.Error("Error while validate login by pattern")
		utils.RespondError(writer, errors.New("bad validate login"), http.StatusBadRequest)

		return
	}

	logging.Logger.Debug("Validate password stage starting....")
	cryptPsq, err := validators.ValidateAndCryptPsw(user.Password, validators.PasswordPhrase)
	if err != nil {
		logging.Logger.Error(err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}
	logging.Logger.Debug("Validate stage finished")
	//set psw to encrypted version
	user.Password = cryptPsq

	logging.Logger.Info("Registering new user in database.")
	if _, err = s.repos.Users.AddNewUser(request.Context(), user); err != nil {
		logging.Logger.Error("Error register new local user. Reason", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Making response for register request")

	utils.Respond(writer, map[string]interface{}{"res": "ok"}, http.StatusOK)
}

func (s *Server) loginRequest(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("loginRequest handler")
	logging.Logger.Info(request.RemoteAddr, " login...")

	user := domain.User{}
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		logging.Logger.Error("request decoding failed. Reason: ", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)
		return
	}

	logging.Logger.Info("Trying authentication user...")
	if user, err = s.repos.Users.FindUser(request.Context(), user.Email, user.Password); err != nil {
		logging.Logger.Error("User doesn't exists or something else. Reason: ", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)
		return
	}

	logging.Logger.Info("Trying add new session by user...")
	session := domain.Session{User: user, Token: uuid.NewString(), SessionTime: time.Now()}
	token, err := s.repos.Sessions.AddSession(request.Context(), session)
	if err != nil {
		logging.Logger.Error("Add new session was failed. Reason: ", err.Error())
		utils.RespondError(writer, err, http.StatusInternalServerError)

		return
	}

	logging.Logger.Debug("Successfully DB operations.")
	logging.Logger.Info("Making response for login request")

	utils.Respond(writer, map[string]interface{}{"token": token}, http.StatusOK)
}

func (s *Server) logoutRequest(writer http.ResponseWriter, request *http.Request) {
	logging.Logger.Debug("logoutRequest handler")
	logging.Logger.Info("Logout user session...")

	token, ok := request.Context().Value("token").(string)
	if !ok {
		logging.Logger.Warning("No token data. This request unauthorized")
		utils.RespondError(writer, errors.New("no token"), http.StatusUnauthorized)

		return
	}

	logging.Logger.Trace(token)

	// delete session from db
	logging.Logger.Info("Deleting session by token")
	err := s.repos.Sessions.DeleteSessionByToken(request.Context(), token)
	if err != nil {
		logging.Logger.Error("Failed deleting session. Reason: ", err.Error())
		utils.RespondError(writer, err, http.StatusBadRequest)

		return
	}

	logging.Logger.Info("Making response for logout request")

	utils.Respond(writer, map[string]interface{}{"res": "ok"}, http.StatusOK)
}
