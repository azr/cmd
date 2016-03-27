// VarHandler wraps variing http handler funcs
//
// Given a func F :
//  func F([ctx context.Context], x X, y Y, z *Z, zz z.Z) {...}
//  // and funcs
//  func HTTPX  (w http.ResponseWriter, r *http.Request) X          {...}
//  func HTTPY  (w http.ResponseWriter, r *http.Request) (Y, error) {...}
//  func HTTPZ  (w http.ResponseWriter, r *http.Request) *Z         {...}
//  func z.HTTPZ(w http.ResponseWriter, r *http.Request) z.Z        {...}
//
// VarHandler will generate an http handler :
//   func FVarHandler(w http.ResponseWriter, r *http.Request) {
//       x := HTTPX(w, r)
//       y, err := HTTPY(w, r)
//       if err != nil {
//       	type HttpErrorInterface interface {
//       		HttpError() (error string, code int)
//       	}
//       	type SelfHttpErrorInterface interface {
//       		HttpError(w http.ResponseWriter)
//       	}
//       	switch t := err.(type) {
//       	default:
//       		w.WriteHeader(http.StatusBadRequest)
//       	case HttpErrorInterface:
//       		error, code := t.HttpError()
//       		t.Error(w, error, code)
//       	case HttpErrorInterface:
//       		t.HttpError(w)
//       	}
//       	return
//       }
//       z := HTTPZ(w, r)
//       zz := z.HTTPZ(w, r)
//       F()
//
//
//   }
//
//
package main // import "github.com/azr/generators/varhandler"

func main() {

}
