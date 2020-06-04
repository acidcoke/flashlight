package middleware


import (
	"github.com/gorilla/sessions"
	"net/http"
)

var store *sessions.CookieStore

func Authenticate(h http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "session")

		// check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(writer, request, "/login", http.StatusFound)
		} else {
			h(writer, request)
		}
	}
}