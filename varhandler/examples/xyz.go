package main

import "net/http"

func HTTPX(r *http.Request) (X, error)  { return X{}, nil }
func HTTPY(r *http.Request) (Y, error)  { return Y{}, nil }
func HTTPZ(r *http.Request) (*Z, error) { return &Z{}, nil }

type X struct {
}
type Y struct {
}
type Z struct {
}
