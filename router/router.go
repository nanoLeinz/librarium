package router

import (
	"net/http"

	"github.com/nanoLeinz/librarium/controller"
	m "github.com/nanoLeinz/librarium/middleware"
)

func NewRouter(member *controller.MemberController,
	auth *controller.AuthController,
	author *controller.AuthorController,
	book *controller.BookController,
	copy *controller.BookCopyController,
	loan *controller.LoanController,
	reservation *controller.ReservationController,
) *http.ServeMux {

	subroute := http.NewServeMux()

	//member
	subroute.Handle("GET /me", m.GenerateTraceID(http.HandlerFunc(member.Profile)))
	subroute.Handle("DELETE /me", m.GenerateTraceID(http.HandlerFunc(member.DeleteProfile)))
	subroute.Handle("PATCH /me", m.GenerateTraceID(http.HandlerFunc(member.UpdateMember)))

	//author
	subroute.Handle("POST /author", m.GenerateTraceID(http.HandlerFunc(author.CreateAuthor)))
	subroute.Handle("GET /author/{id}", m.GenerateTraceID(http.HandlerFunc(author.GetByID)))
	subroute.Handle("GET /author", m.GenerateTraceID(m.Paginator(http.HandlerFunc(author.GetAllAuthor))))
	subroute.Handle("DELETE /author/{id}", m.GenerateTraceID(http.HandlerFunc(author.DeleteByID)))
	subroute.Handle("PATCH /author/{id}", m.GenerateTraceID(http.HandlerFunc(author.UpdateAuthor)))
	subroute.Handle("GET /author/{id}/books", m.GenerateTraceID(http.HandlerFunc(author.GetAuthorsBook)))

	//book
	subroute.Handle("POST /book", m.GenerateTraceID(http.HandlerFunc(book.CreateBook)))
	subroute.Handle("DELETE /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.DeleteBook)))
	subroute.Handle("PATCH /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.UpdateBook)))
	subroute.Handle("GET /book/{id}", m.GenerateTraceID(http.HandlerFunc(book.GetBook)))
	subroute.Handle("GET /book", m.GenerateTraceID(m.Paginator(http.HandlerFunc(book.GetAll))))
	subroute.Handle("GET /book/search", m.GenerateTraceID(m.Paginator(http.HandlerFunc(book.GetBookByTitle))))

	//book copy
	subroute.Handle("POST /book/{bookID}/copies", m.GenerateTraceID(http.HandlerFunc(copy.CreateCopies)))
	subroute.Handle("DELETE /book/{bookID}/copies/{copyID}", m.GenerateTraceID(http.HandlerFunc(copy.DeleteCopy)))
	subroute.Handle("PATCH /book/{bookID}/copies/{copyID}", m.GenerateTraceID(http.HandlerFunc(copy.UpdateStatus)))
	subroute.Handle("GET /book/{bookID}/copies", m.GenerateTraceID(m.Paginator(http.HandlerFunc(copy.GetCopyByCondition))))
	subroute.Handle("GET /book/copies", m.GenerateTraceID(m.Paginator(http.HandlerFunc(copy.GetAll))))
	subroute.Handle("GET /book/{bookID}/copies/{copyID}", m.GenerateTraceID(http.HandlerFunc(copy.GetCopy)))

	//loan
	subroute.Handle("POST /loans", m.GenerateTraceID(http.HandlerFunc(loan.CreateLoan)))
	subroute.Handle("DELETE /loans/{id}", m.GenerateTraceID(http.HandlerFunc(loan.DeleteLoan)))
	subroute.Handle("PATCH /loans/{id}", m.GenerateTraceID(http.HandlerFunc(loan.UpdateLoan)))
	subroute.Handle("GET /loans/{id}", m.GenerateTraceID(http.HandlerFunc(loan.GetLoanByID)))
	subroute.Handle("GET /loans", m.GenerateTraceID(m.Paginator(http.HandlerFunc(loan.GetAllLoan))))

	//reservation
	subroute.Handle("POST /reservation", m.GenerateTraceID(http.HandlerFunc(reservation.CreateReservation)))
	subroute.Handle("DELETE /reservation/{id}", m.GenerateTraceID(http.HandlerFunc(reservation.DeleteReservation)))
	subroute.Handle("PATCH /reservation/{id}", m.GenerateTraceID(http.HandlerFunc(reservation.UpdateReservation)))
	subroute.Handle("GET /reservation/{id}", m.GenerateTraceID(http.HandlerFunc(reservation.GetReservationByID)))
	subroute.Handle("GET /reservation", m.GenerateTraceID(m.Paginator(http.HandlerFunc(reservation.GetAllReservation))))

	//v1 api
	mainroute := http.NewServeMux()
	mainroute.Handle("/api/v1/", m.ExtendContext(m.ValidateJWT(http.StripPrefix("/api/v1", subroute))))

	//auth
	mainroute.HandleFunc("POST /api/v1/members", auth.Register)
	mainroute.HandleFunc("POST /api/v1/login", auth.Login)

	return mainroute

}
