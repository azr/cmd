//go:generate varhandler -func Simple,Import,Status,Response,ResponseStatus,CreateUser,GetUser,UpdateUser,DeleteUser
package main

import "net/http"

func main() {
	http.ListenAndServe(":8080", nil)
}
