package z

import "net/http"

type Z struct {
}

func HTTPZ(r *http.Request) (Z, error) {
	return Z{}, nil
}
