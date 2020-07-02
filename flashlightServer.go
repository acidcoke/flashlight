package main

import (
	"flashlight/app/controller"
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

	http.HandleFunc("/add-user", controller.AddUser)
	http.HandleFunc("/authenticate-user", controller.AuthenticateUser)
	http.HandleFunc("/add-post", controller.AddPost)
	http.HandleFunc("/like-post", controller.LikePost)
	http.HandleFunc("/delete-post", controller.DeletePost)
	http.HandleFunc("/add-comment", controller.AddComment)
	http.HandleFunc("/delete-comment", controller.DeleteComment)
	http.HandleFunc("/logout", controller.Logout)

	// fileserver path mapping
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))

	http.ListenAndServe(":8080", nil)
}
