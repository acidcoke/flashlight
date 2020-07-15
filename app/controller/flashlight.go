package controller

import (
	"flashlight/app/model"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/pborman/uuid"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

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

func AddPost(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		request.ParseMultipartForm(32 << 20)
		// "datei" ist das Attribut name des Html-Input-Tags
		file, _, err := request.FormFile("datei")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		nameOfImage := uuid.New()
		imagePath := "../data/pictures/" + nameOfImage + ".jpg"
		//Bild wird in Ordner "test" gespeichert. Ordner muss zuvor angelegt werden!
		f, err := os.OpenFile("./data/pictures/"+nameOfImage+".jpg", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		session, _ := store.Get(request, "session")
		username := session.Values["username"].(string)
		caption := request.FormValue("caption")
		timestamp := time.Now().Format("02.01.2006 - 15:04")

		flashlight := model.Flashlight{
			Type:       "Flashlight",
			FilePath:   imagePath,
			Author:     username,
			Timestamp:  timestamp,
			LikeAmount: 0,
			Caption:    caption,
		}
		flashlight.Add()
		http.Redirect(w, request, "/home", http.StatusFound)
	}
}

func DeletePost(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		id := request.FormValue("id")

		_ = model.Delete(id)
		_ = model.DeleteCommentByFlashlightId(id)
		http.Redirect(w, request, "/mypictures", http.StatusFound)
	}
}

func LikePost(w http.ResponseWriter, request *http.Request) {

	id := request.FormValue("fid")
	session, _ := store.Get(request, "session")
	username := session.Values["username"].(string)
	var like model.Like


		like = model.Like{
			Type:         "Like",
			FlashlightId: id,
			Username:     username,
		}
		_ = model.AddLike(like)

	_, _ = model.CountLikes(id)
	http.Redirect(w, request, "/home", http.StatusFound)
}

func DislikePost(w http.ResponseWriter, request *http.Request) {

	id := request.FormValue("fid")
	session, _ := store.Get(request, "session")
	username := session.Values["username"].(string)
	var like model.Like


		like, _ = model.GetLike(username, id)
		_ = model.DeleteLike(like.ID)

	_, _ = model.CountLikes(id)
	http.Redirect(w, request, "/home", http.StatusFound)
}

func AddComment(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		session, _ := store.Get(request, "session")
		username := session.Values["username"].(string)
		commentText := request.FormValue("comment")
		id := request.FormValue("fid")

		comment := model.Comment{
			Type:         "Comment",
			Author:       username,
			Text:         commentText,
			FlashlightId: id,
		}
		model.AddComment(comment)
		http.Redirect(w, request, "/home", http.StatusFound)
	}
}

func DeleteComment(w http.ResponseWriter, request *http.Request) {
}

func Index(writer http.ResponseWriter, request *http.Request) {
	flashlights, _ := model.GetAllFlashlights()
	for _, flashlight := range flashlights {
		flashlight.Comments, _ = flashlight.GetComments()
	}
	data := struct {
		Flashlights *[]model.Flashlight
	}{
		&flashlights,
	}
	tmpl.ExecuteTemplate(writer, "flashlight.tmpl", data)
}

func Login(writer http.ResponseWriter, request *http.Request) {
	tmpl.ExecuteTemplate(writer, "login.tmpl", 11)
}

func Home(writer http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "session")
	username := session.Values["username"].(string)
	flashlights, _ := model.GetAllFlashlights()
	for index := range flashlights{
		flashlights[index].Comments, _ =flashlights[index].GetComments()
		like, _ := model.GetLike(username, flashlights[index].ID)
		if like.Username != "" {
			flashlights[index].IsLiked=1
		} else {
			flashlights[index].IsLiked=0
		}
	}
	data := struct {
		Flashlights *[]model.Flashlight
		Username    string
	}{
		&flashlights,
		username,
	}
	tmpl.ExecuteTemplate(writer, "home.tmpl", data)
}

func MyPictures(writer http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "session")
	username := session.Values["username"].(string)
	flashlights, _ := model.GetFlashlightsByUser(username)
	data := struct {
		Flashlights *[]model.Flashlight
		Username    string
	}{
		&flashlights,
		username,
	}
	tmpl.ExecuteTemplate(writer, "mypictures.tmpl", data)
}

func Registration(writer http.ResponseWriter, request *http.Request) {
	tmpl.ExecuteTemplate(writer, "registration.tmpl", 11)
}

func Upload(writer http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "session")
	username := session.Values["username"].(string)
	data := struct {
		Username string
	}{
		username,
	}
	tmpl.ExecuteTemplate(writer, "upload.tmpl", data)
}
