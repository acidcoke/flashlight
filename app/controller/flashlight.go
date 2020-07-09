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
		timestamp := time.Now()

		flashlight := model.Flashlight{
			Type:       "Flashlight",
			FilePath:   imagePath,
			Author:     username,
			Timestamp:  timestamp,
			LikeAmount: 1,
			Caption:    caption,
		}
		flashlight.Add()
		Home(w, request)
	}
}

func DeletePost(w http.ResponseWriter, request *http.Request) {

}

func LikePost(w http.ResponseWriter, request *http.Request) {
}

func AddComment(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		session, _ := store.Get(request, "session")
		username := session.Values["username"].(string)
		commentText := request.FormValue("comment")
		id := request.FormValue("test")

		comment := model.Comment{
			Type:         "Comment",
			Author:       username,
			Text:         commentText,
			FlashlightId: id,
		}
		model.AddComment(comment)
		Home(w, request)
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
