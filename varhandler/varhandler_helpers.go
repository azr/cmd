// Code copyied from github.com/azr/generators/varhandler/varhandler_helpers.go; DO NOT EDIT
package main

import "net/http"

func HandleHttpErrorWithDefaultStatus(w http.ResponseWriter, r *http.Request, status int, err error) {
	type HttpError interface {
		HttpError() (error string, code int)
	}
	type SelfHttpError interface {
		HttpError(w http.ResponseWriter)
	}
	switch t := err.(type) {
	default:
		w.WriteHeader(status)
	case HttpError:
		err, code := t.HttpError()
		http.Error(w, err, code)
	case http.Handler:
		t.ServeHTTP(w, r)
	case SelfHttpError:
		t.HttpError(w)
	}
}

func HandleHttpResponse(w http.ResponseWriter, r *http.Request, resp interface{}) {
	type Byter interface {
		Bytes() []byte
	}
	type Stringer interface {
		String() string
	}
	switch t := resp.(type) {
	default:
		// I don't know that type !
	case http.Handler:
		t.ServeHTTP(w, r) // resp knows how to handle itself
	case Byter:
		w.Write(t.Bytes())
	case Stringer:
		w.Write([]byte(t.String()))
	case []byte:
		w.Write(t)
	}
}
