package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Kran001/basic-auth/pkg/logging"
)

// Respond - create HTTP answer with data field and stCode for status
func Respond(w http.ResponseWriter, data map[string]interface{}, stCode int) {
	w.Header().Add("Content-Type", "application/json")

	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods-Type", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.WriteHeader(stCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		if data != nil {
			logging.Logger.Warning("RESPONSE ERROR: ", err, len(data))
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]interface{}{"code": http.StatusBadRequest, "message": "no data"}})
		if err != nil {
			logging.Logger.Warning("ERROR spare resp")
		}
	}
}

// RespondError - Sending error response
func RespondError(writer http.ResponseWriter, err error, status int) {
	nErr := MakeErr(err.Error())
	response := map[string]interface{}{"error": map[string]interface{}{"code": nErr.Code, "msg": nErr.Msg}}
	Respond(writer, response, status)
}

// CheckToken - Checking token wrap.
func CheckToken(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		token, ok := request.Header["Token"]
		if !ok || len(token) < 1 {
			logging.Logger.Error("No Auth token")
			Respond(writer,
				map[string]interface{}{
					"error": map[string]interface{}{
						"code": http.StatusUnauthorized, "msg": "no token",
					},
				}, http.StatusUnauthorized)
		}

		request = request.WithContext(context.WithValue(request.Context(), "token", token[0]))

		h.ServeHTTP(writer, request)
	}
}
