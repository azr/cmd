//go:generate varhandler -func CreateUser,GetUser,UpdateUser,DeleteUser -output user_handlers_generated.go
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/user/create", CreateUserHandler)
	http.HandleFunc("/user/get", GetUserHandler)
	http.HandleFunc("/user/update", UpdateUserHandler)
	http.HandleFunc("/user/delete", DeleteUserHandler)
}

///////
/// Types
///////

// User

type User struct {
	Id   UserID
	Name string
}

func HTTPUser(r *http.Request) (u User, err error) {
	// if request encoding is json :
	err = json.NewDecoder(r.Body).Decode(&u)
	if u.Name == "" {
		return u, errors.New("EmptyName")
	}
	return
}

func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if response encoding has to be json :
	json.NewEncoder(w).Encode(u)
}

//UserID
type UserID string

func HTTPUserID(r *http.Request) (uid UserID, err error) {
	uid = UserID(r.URL.Query().Get("user_id"))
	if uid == "" {
		return uid, errors.New("Please provide a user id")
	}
	log.Printf("uid: %s", uid)
	return
}

///////
/// Handlers
///////

//create

func CreateUser(user User) (status int, err error) {
	//save user into database
	return http.StatusCreated, err
	//if the error is not nil
	//a status internal server
	//error will be returned by default
	//otherwise, everything is fine.
}

//get

func GetUser(id UserID) (resp http.Handler, status int, err error) {
	if id == "404" { // check case
		return nil, http.StatusNotFound, nil
	}
	user := User{
		Id: id,
	}
	//err := db.GetUser(user)
	return user, http.StatusOK, err
}

//update

func UpdateUser(id UserID, user User) (status int, err error) {
	//user might have to be
	//a UserUpdateRequest type
	//that only takes into account modifiable fields

	if id == "404" { // check case
		return http.StatusNotFound, nil
	}
	user.Id = id
	//err := db.GetUser(user)
	if err != nil {
		return
	}
	//do stuff

	return http.StatusOK, nil
}

//delete

func DeleteUser(id UserID) (status int, err error) {
	if id == "404" { // check case
		return http.StatusNotFound, nil
	}
	//db.DeleteUser(id)
	if err != nil {
		return
	}

	return http.StatusNoContent, nil
}
