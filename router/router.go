package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
	"github.com/nanoLeinz/librarium/middleware"
)

func NewRouter(member *controller.MemberController, auth *controller.AuthController, author *controller.AuthorController, book *controller.BookController) *http.ServeMux {

	subroute := http.NewServeMux()

	//auth
	subroute.HandleFunc("POST /members", auth.Register)
	subroute.HandleFunc("POST /login", auth.Login)

	//member
	subroute.Handle("GET /me", middleware.ValidateJWT(http.HandlerFunc(member.Profile)))
	subroute.Handle("DELETE /me", middleware.ValidateJWT(http.HandlerFunc(member.DeleteProfile)))
	subroute.Handle("PATCH /me", middleware.ValidateJWT(http.HandlerFunc(member.UpdateMember)))

	//book
	subroute.Handle("POST /book", middleware.ValidateJWT(http.HandlerFunc(book.CreateBook)))

	//author
	subroute.Handle("POST /author", middleware.ValidateJWT(http.HandlerFunc(author.CreateAuthor)))

	//v1 api
	mainroute := http.NewServeMux()
	mainroute.Handle("/api/v1/", http.StripPrefix("/api/v1", subroute))

	return mainroute

}
