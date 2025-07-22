package helper

import (
	"encoding/json"
	"net/http"

	"github.com/nanoLeinz/librarium/model/dto"
)

func ResponseJSON(w http.ResponseWriter, data *dto.WebResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.Code)
	_ = json.NewEncoder(w).Encode(data)
}
