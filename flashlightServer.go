package main

import (
	"flashlight/app/controller"
	"net/http"
)

func main() {

	// function mapping
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/home", controller.Auth(controller.Home))
	http.HandleFunc("/mypictures", controller.Auth(controller.MyPictures))
	http.HandleFunc("/upload", controller.Auth(controller.Upload))
	http.HandleFunc("/registration", controller.Registration)

	http.HandleFunc("/add-user", controller.AddUser)
	http.HandleFunc("/authenticate-user", controller.AuthenticateUser)
	http.HandleFunc("/add-post", controller.Auth(controller.AddPost))
	http.HandleFunc("/like-post", controller.Auth(controller.LikePost))
	http.HandleFunc("/delete-post/", controller.Auth(controller.DeletePost))
	http.HandleFunc("/add-comment", controller.Auth(controller.AddComment))
	http.HandleFunc("/delete-comment", controller.Auth(controller.DeleteComment))
	http.HandleFunc("/logout", controller.Logout)

	// fileserver path mapping
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data"))))

	http.ListenAndServe(":80", nil)
}
