package controller

import "net/http"

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := model.User{}
	user.Username = username
	user.Password = password

	user.Add()

	tmpl.ExecuteTemplate(w, "login.tmpl", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
}
