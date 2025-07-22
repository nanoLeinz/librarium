package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
)

func NewRouter(member *controller.MemberController, auth *controller.AuthController) *http.ServeMux {

	subroute := http.NewServeMux()

	subroute.HandleFunc("POST /members", auth.Register)
	subroute.HandleFunc("POST /login", auth.Login)

	//v1 api
	mainroute := http.NewServeMux()
	mainroute.Handle("/api/v1/", http.StripPrefix("/api/v1", subroute))

	return mainroute

}
