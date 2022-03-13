package utils

import (
	"encoding/json"
	log "github.com/Kran001/basic-auth/pkg/logging"
	"net/http"
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
			log.Logger.Warning("RESPONSE ERROR: ", err, len(data))
		}
		err = json.NewEncoder(w).Encode(map[string]interface{}{"error": map[string]interface{}{"code": http.StatusBadRequest, "message": "no data"}})
		if err != nil {
			log.Logger.Warning("ERROR spare resp")
		}
	}
}
