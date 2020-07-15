package model

import (
	"encoding/json"
	"fmt"
	"github.com/leesper/couchdb-golang"
	"runtime"
	"strings"
)

type Flashlight struct {
	ID         string    `json:"_id"`
	Rev        string    `json:"_rev"`
	Type       string    `json:"type"`
	FilePath   string    `json:"file_path"`
	Author     string    `json:"author"`
	Timestamp  string `json:"timestamp"`
	LikeAmount int       `json:"like_amount"`
	Caption    string    `json:"caption"`
	Comments   []Comment
	IsLiked		int
}

type Comment struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev"`
	Type         string `json:"type"`
	Author       string `json:"author"`
	Text         string `json:"text"`
	FlashlightId string `json:"flashlight_id"`
}

type Like struct {
	ID         		string  `json:"_id"`
	Rev        		string  `json:"_rev"`
	Type       		string  `json:"type"`
	FlashlightId 	string	`json:"flashlight_id"`
	Username 		string 	`json:"username"`
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
		printError(err)
	}
	return err
}

func AddComment(c Comment) (err error) {
	comment, _ := comment2Map(c)
	delete(comment, "_id")
	delete(comment, "_rev")
	_, _, err = flashlightDb.Save(comment, nil)

	if err != nil {
		printError(err)
	}
	return err
}

func AddLike(like Like) (err error) {
	dbLike, _ := GetLike(like.Username, like.FlashlightId)
	if dbLike.Username == "" {
		likeMap, _ := like2Map(like)
		delete(likeMap, "_id")
		delete(likeMap, "_rev")

		_, _, err = flashlightDb.Save(likeMap, nil)
		if err != nil {
			printError(err)
		}
	} else {
		_ = DeleteLike(dbLike.ID)
	}
	return err
}

func GetLike(username string, flashlightId string) (like Like, err error) {
	query := `
	{
		"selector": {
			"type": "Like",
			"flashlight_id": "%s",
			"username": "%s"
		}
	 }`
	likeMap, err := flashlightDb.QueryJSON(fmt.Sprintf(query, flashlightId, username))
	if len(likeMap) > 0 {
		like, _ = map2Like(likeMap[0])
	} else {
		like.Username = ""
	}
	if err != nil {
		printError(err)
	}
	return like, err
}

func DeleteLike(id string) (err error) {
	err = flashlightDb.Delete(id)
	if err != nil {
		printError(err)
	}
	return err
}

func DeleteLikeByUserId(id string) (err error) {
	query := `
	{
		"selector": {
			"type": "Like",
			"user_id": "%s"
		}
	 }`
	likes, err := flashlightDb.QueryJSON(fmt.Sprintf(query, id))
	if err != nil {
		printError(err)
		return err
	}
	for index := range likes {
		err = DeleteLike(likes[index]["_id"].(string))
	}
	if err != nil {
		printError(err)
		
	}
	return err
}

/*func DeleteLikeByFlashlightId(id string) error {

}*/

func CountLikes(flashlightId string) (likeCount int, err error) {
	query := `
	{
		"selector": {
			"type": "Like",
			"flashlight_id": "%s"
		}
	 }`
	likeMap, err := flashlightDb.QueryJSON(fmt.Sprintf(query, flashlightId))

	if err != nil {
		printError(err)
	}
	flashlight, _ :=GetFlashlight(flashlightId)
	flashlight.LikeAmount=len(likeMap)
	var x []map[string]interface{}
	flashlightMap, _ := flashlight2Map(flashlight)
	x = append(x, flashlightMap)
	flashlightDb.Update(x, nil)

	return len(likeMap), err
}

/*func UpdateLikeCount(flashlightId string, likeCount int) (err error) {
	flashlightDb.Get(flashlightId)
}*/

func DeleteComment(id string) (err error) {
	err = flashlightDb.Delete(id)
	if err != nil {
		printError(err)
	}
	return err
}

func DeleteCommentByFlashlightId(id string) (err error) {
	query := `
	{
		"selector": {
			"type": "Comment",
			"flashlight_id": "%s"
		}
	 }`

	commentMaps, err := flashlightDb.QueryJSON(fmt.Sprintf(query, id))
	if err != nil {
		printError(err)
	}
	for index := range commentMaps {
		_ = flashlightDb.Delete(commentMaps[index]["_id"].(string))
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

func Delete(id string) error {
	err := flashlightDb.Delete(id)
	if err != nil {
		printError(err)
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

func like2Map(l Like) (like map[string]interface{}, err error) {
	flashlightJSON, err := json.Marshal(l)
	json.Unmarshal(flashlightJSON, &like)

	return like, err
}

func map2Like(like map[string]interface{}) (l Like, err error) {
	flashlightJSON, err := json.Marshal(like)
	json.Unmarshal(flashlightJSON, &l)

	return l, err
}

func printError(err error) {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		result := strings.Split(details.Name(), ".")
		fmt.Printf("[%s] error: %s\n", result[1], err)
	}
}