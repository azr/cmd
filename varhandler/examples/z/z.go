package z

import "net/http"

type Z struct {
}

func HTTPZ(w http.ResponseWriter, r *http.Request) (Z, error) {
	return Z{}, nil
}
