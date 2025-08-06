package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
	"github.com/nanoLeinz/librarium/middleware"
)

func NewRouter(member *controller.MemberController, auth *controller.AuthController) *http.ServeMux {

	subroute := http.NewServeMux()

	subroute.HandleFunc("POST /members", auth.Register)
	subroute.HandleFunc("POST /login", auth.Login)

	subroute.Handle("GET /me", middleware.ValidateJWT(http.HandlerFunc(member.Profile)))
	subroute.Handle("DELETE /me", middleware.ValidateJWT(http.HandlerFunc(member.DeleteProfile)))
	subroute.Handle("PATCH /me", middleware.ValidateJWT(http.HandlerFunc(member.UpdateMember)))

	//v1 api
	mainroute := http.NewServeMux()
	mainroute.Handle("/api/v1/", http.StripPrefix("/api/v1", subroute))

	return mainroute

}
