package main

import (
	"flash2/app/controller"
	"net/http"
)

func main() {

	// function mapping
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/home", controller.Home)
	http.HandleFunc("/mypictures", controller.MyPictures)
	http.HandleFunc("/upload", controller.Upload)
	http.HandleFunc("/registration", controller.Registration)

	// fileserver path mapping
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))

	http.ListenAndServe(":8080", nil)
}
