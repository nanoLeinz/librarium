package controller

import (
	"encoding/json"
	"net/http"
)

func CreateMember(w http.ResponseWriter, r *http.Request) {

	var req 
	
	json.NewDecoder(r.Body).Decode(req)

}
