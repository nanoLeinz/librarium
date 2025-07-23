package helper

import (
	"encoding/json"
	"net/http"

	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"
)

func ResponseJSON(w http.ResponseWriter, data *dto.WebResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.Code)
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.WithFields(
			log.Fields{
				"data": data,
			},
		).Error("error while encoding the data")
	}

}
