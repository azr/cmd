package main

import (
	"net/http"
	z "z"

	"golang.org/x/net/context"
)

func F(ctx context.Context, x X, y Y, z *Z, zz z.Z) {
}

func HTTPX(w http.ResponseWriter, r *http.Request) X          { return X{} }
func HTTPY(w http.ResponseWriter, r *http.Request) (Y, error) { return Y{}, nil }
func HTTPZ(w http.ResponseWriter, r *http.Request) *Z         { return &Z{} }

type X struct {
}
type Y struct {
}
type Z struct {
}
