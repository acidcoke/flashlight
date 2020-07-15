package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetUserByUsername(username string) (user User, err error) {
	user = User{}

	query := `
	{
		"selector": {
			"type": "User",
			"username": "%s"
		}
	}`

	u, err := flashlightDb.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return user, err
	}
	if len(u) != 0 {
		user, err = map2User(u[0])
		if err != nil {
			return User{}, err
		}
		return user, nil
	}
	return user, fmt.Errorf("Fehler")

}

func (user User) Add() (err error) {
	// Check wether username already exists
	// Todo...
	existinguser, _ := GetUserByUsername(user.Username)
	if existinguser.Id == "" {
		// Hash password
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

		user.Password = b64HashedPwd
		user.Type = "User"

		// Convert user struct to map[string]interface as required by Save() method
		u, err := user2Map(user)

		// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
		delete(u, "_id")
		delete(u, "_rev")

		// Add user to DB
		_, _, err = flashlightDb.Save(u, nil)

		if err != nil {
			fmt.Printf("[Add] error: %s", err)
		}
	} else {
		err := fmt.Errorf("user already exists")
		fmt.Println(err.Error())
	}
	return err
}

// ---------------------------------------------------------------------------
// Internal helper functions
// ---------------------------------------------------------------------------

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func user2Map(u User) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2User(user map[string]interface{}) (u User, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
