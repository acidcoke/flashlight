package controller

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"

	"crypto/rand"
	_ "encoding/base64"
)

// struct for template execution
type PageData struct {
	PageTitle string
}

var tmpl *template.Template

var store *sessions.CookieStore

//automatically(once) called during package startup
func init() {
	tmpl = template.Must(template.ParseGlob("template/*.tmpl"))

	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewCookieStore(key)

}

// AddUser controller

func AddPost(w http.ResponseWriter, r *http.Request) {
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
}

func LikePost(w http.ResponseWriter, r *http.Request) {
}

func AddComment(w http.ResponseWriter, r *http.Request) {
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
}

func Index(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "flashlight.tmpl", PageData{"G'day m8"})
}

func Login(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "login.tmpl", 11)
}

func Home(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "home.tmpl", 11)
}

func MyPictures(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "mypictures.tmpl", 11)
}

func Registration(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "registration.tmpl", 11)
}

func Upload(writer http.ResponseWriter, reader *http.Request) {
	tmpl.ExecuteTemplate(writer, "upload.tmpl", 11)
}
