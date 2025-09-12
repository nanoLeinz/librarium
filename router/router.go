package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
	m "github.com/nanoLeinz/librarium/middleware"
)

func NewRouter(member *controller.MemberController, auth *controller.AuthController, author *controller.AuthorController, book *controller.BookController) *http.ServeMux {

	subroute := http.NewServeMux()

	//member
	subroute.Handle("GET /me", m.GenerateTraceID(http.HandlerFunc(member.Profile)))
	subroute.Handle("DELETE /me", m.GenerateTraceID(http.HandlerFunc(member.DeleteProfile)))
	subroute.Handle("PATCH /me", m.GenerateTraceID(http.HandlerFunc(member.UpdateMember)))

	//book
	subroute.Handle("POST /book", m.GenerateTraceID(http.HandlerFunc(book.CreateBook)))
	subroute.Handle("DELETE /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.DeleteBook)))
	subroute.Handle("PATCH /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.UpdateBook)))
	subroute.Handle("GET /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.GetBook)))
	subroute.Handle("GET /book", m.GenerateTraceID(m.Paginator(http.HandlerFunc(book.GetAll))))
	subroute.Handle("GET /book/search", m.GenerateTraceID(http.HandlerFunc(book.GetBookByTitle)))

	//author
	subroute.Handle("POST /author", m.GenerateTraceID(http.HandlerFunc(author.CreateAuthor)))

	//v1 api
	mainroute := http.NewServeMux()
	mainroute.Handle("/api/v1/", m.ExtendContext(m.ValidateJWT(http.StripPrefix("/api/v1", subroute))))

	//auth
	mainroute.HandleFunc("POST /api/v1/members", auth.Register)
	mainroute.HandleFunc("POST /api/v1/login", auth.Login)

	return mainroute

}
