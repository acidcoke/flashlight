package controller

import (
	"encoding/base64"
	"flashlight/app/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func AuthenticateUser(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	user, _ := model.GetUserByUsername(username)
	// decode base64 String to []byte
	passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
	err := bcrypt.CompareHashAndPassword(passwordDB, []byte(password))

	if err == nil {
		session, _ := store.Get(request, "session")

		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(request, writer)
		http.Redirect(writer, request, "/home", http.StatusFound)
	} else {
		tmpl.ExecuteTemplate(writer, "login.tmpl", "nil")
	}
}

// Logout controller
func Logout(writer http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "session")

	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Save(request, writer)

	tmpl.ExecuteTemplate(writer, "flashlight.tmpl", nil)
}

func AddUser(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	user := model.User{}
	user.Username = username
	user.Password = password

	user.Add()

	tmpl.ExecuteTemplate(writer, "login.tmpl", nil)
}

// Authentication handler
func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		session, _ := store.Get(request, "session")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(writer, request, "/login", http.StatusFound)
		} else {
			handler(writer, request)
		}
	}
}
