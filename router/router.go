package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
)

func NewRouter(member *controller.MemberController) *http.ServeMux {

	route := http.NewServeMux()

	route.HandleFunc("POST /api/v1/members", member.CreateMember)

	return route

}
