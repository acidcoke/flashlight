package model

import (
	"encoding/json"
	"fmt"
	"github.com/leesper/couchdb-golang"
	"time"
)

type Flashlight struct {
	ID         string    `json:"_id"`
	Rev        string    `json:"_rev"`
	Type       string    `json:"type"`
	FilePath   string    `json:"file_path"`
	Author     string    `json:"author"`
	Timestamp  time.Time `json:"timestamp"`
	LikeAmount int       `json:"like_amount"`
	Caption    string    `json:"caption"`
	Comments   []Comment
}

type Comment struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev"`
	Type         string `json:"type"`
	Author       string `json:"author"`
	Text         string `json:"text"`
	FlashlightId string `json:"flashlight_id"`
}

var flashlightDb *couchdb.Database

func init() {
	var err error
	flashlightDb, err = couchdb.NewDatabase("http://localhost:5984/flashlight")
	if err != nil {
		panic(err)
	}
}

func (f Flashlight) Add() error {
	flashlight, _ := flashlight2Map(f)
	delete(flashlight, "_id")
	delete(flashlight, "_rev")
	_, _, err := flashlightDb.Save(flashlight, nil)
	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}
	return err
}

func AddComment(c Comment) (err error) {
	comment, _ := comment2Map(c)
	delete(comment, "_id")
	delete(comment, "_rev")
	_, _, err = flashlightDb.Save(comment, nil)

	if err != nil {
		fmt.Printf("[AddComment] error: %s", err)
	}
	return err
}

func (c Comment) DeleteComment() (err error) {
	err = flashlightDb.Delete(c.ID)
	if err != nil {
		fmt.Printf("[DeleteComment] error: %s", err)
	}
	return err
}

func (f Flashlight) GetComments() (comments []Comment, err error) {
	query := `
	{
		"selector": {
			"type": "Comment",
			"flashlight_id": "%s"
		}
	 }`

	commentMaps, err := flashlightDb.QueryJSON(fmt.Sprintf(query, f.ID))
	if err != nil {
		return nil, err
	}
	for index := range commentMaps {
		comment, _ := map2Comment(commentMaps[index])
		comments = append(comments, comment)
	}
	return comments, nil
}

/*func (c Comment) EditComment() (err error){

}*/

func GetFlashlight(id string) (flashlight Flashlight, err error) {
	f, err := flashlightDb.Get(id, nil)
	if err != nil {
		return Flashlight{}, err
	}
	flashlight, _ = map2Flashlight(f)
	return flashlight, nil
}

func GetFlashlightsByUser(username string) (flashlights []Flashlight, err error) {
	query := `
	{
		"selector": {
			"type": "Flashlight",
			"author": "%s"
		}
	 }`

	flashlightMaps, err := flashlightDb.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return nil, err
	} else {
		for index := range flashlightMaps {
			flashlight, _ := map2Flashlight(flashlightMaps[index])
			flashlight.Comments, _ = flashlight.GetComments()
			flashlights = append(flashlights, flashlight)
		}
		return flashlights, nil
	}
}

func GetAllFlashlights() (flashlights []Flashlight, err error) {
	flashlightMaps, err := flashlightDb.QueryJSON(`
	{
		"selector": {
			 "type": {
					"$eq": "Flashlight"
			 }
		}
	 }`)
	if err != nil {
		return nil, err
	} else {
		for index := range flashlightMaps {
			flashlight, _ := map2Flashlight(flashlightMaps[index])
			flashlight.Comments, _ = flashlight.GetComments()
			flashlights = append(flashlights, flashlight)
		}
		return flashlights, nil
	}

}

func (f Flashlight) Delete() error {
	err := flashlightDb.Delete(f.ID)
	if err != nil {
		fmt.Printf("[Delete] error: %s", err)
	}
	return err
}

//Internal functions

func flashlight2Map(f Flashlight) (flashlight map[string]interface{}, err error) {
	flashlightJSON, err := json.Marshal(f)
	json.Unmarshal(flashlightJSON, &flashlight)

	return flashlight, err
}

func map2Flashlight(flashlight map[string]interface{}) (f Flashlight, err error) {
	flashlightJSON, err := json.Marshal(flashlight)
	json.Unmarshal(flashlightJSON, &f)

	return f, err
}

func comment2Map(c Comment) (comment map[string]interface{}, err error) {
	flashlightJSON, err := json.Marshal(c)
	json.Unmarshal(flashlightJSON, &comment)

	return comment, err
}

func map2Comment(comment map[string]interface{}) (c Comment, err error) {
	flashlightJSON, err := json.Marshal(comment)
	json.Unmarshal(flashlightJSON, &c)

	return c, err
}
