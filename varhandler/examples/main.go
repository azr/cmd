//go:generate varhandler -func Status,Response,ResponseStatus
package main

import "net/http"

func main() {
	http.ListenAndServe(":8080", nil)
}
