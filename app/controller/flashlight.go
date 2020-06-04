package controller

import (
	"html/template"
	"net/http"
)

// struct for template execution
type PageData struct {
	PageTitle string
}

var tmpl *template.Template

//automatically(once) called during package startup
func init(){
	tmpl = template.Must(template.ParseGlob("template/*.tmpl"))
}

func Index(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "flashlight.tmpl", PageData{"G'day m8"})
}

func Login(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "login.tmpl", 11)
}

func Home(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "home.tmpl", 11)
}

func MyPictures(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "mypictures.tmpl", 11)
}

func Registration(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "registration.tmpl", 11)
}

func Upload(writer http.ResponseWriter, reader *http.Request){
	tmpl.ExecuteTemplate(writer, "upload.tmpl", 11)
}