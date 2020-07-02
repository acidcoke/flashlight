package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/leesper/couchdb-golang"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var flashlightDb *couchdb.Database

func init() {
	var err error
	flashlightDb, err = couchdb.NewDatabase("http://localhost:5984/flashlight")
	if err != nil {
		panic(err)
	}
}

func (user User) Add() (err error) {
	// Check wether username already exists
	// Todo...

	// Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

	user.Password = b64HashedPwd
	user.Type = "User"

	// Convert Todo struct to map[string]interface as required by Save() method
	u, err := user2Map(user)

	// Delete _id and _rev from map, otherwise DB access will be denied (unauthorized)
	delete(u, "_id")
	delete(u, "_rev")

	// Add todo to DB
	_, _, err = flashlightDb.Save(u, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
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
