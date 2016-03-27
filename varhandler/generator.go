// VarHandler genrate wrappers for variing http handler funcs
//
// Given a pkg :
// 	//go:generate varhandler
//  package server
//
//  func F([ctx context.Context], x X, y Y, z *Z, zz z.Z) ([resp interface{},] [status int,] err error) {...}
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
//       	HandleHttpErrorWithDefaultStatus(http.StatusBadRequest, err)
//       	return
//       }
//       z := HTTPZ(w, r)
//       zz := z.HTTPZ(w, r)
//       resp, status, err := F([context.Background()], x, y, z, zz)
//       if err != nil {
//       	HandleHttpErrorWithDefaultStatus(http.InternalServerError, err)
//       	return
//       }
//       if status != 0 { // code generated if status is returned
//       	w.WriteHeader(status)
//       }
//       if resp != nil { // code generated if resp object is returned
//       	HandleHttpResponse()
//       }
//   }
//
//   func HandleHttpErrorWithDefaultStatus(status int, err error) {
//       type HttpError interface {
//       	HttpError() (error string, code int)
//       }
//       type SelfHttpError interface {
//       	HttpError(w http.ResponseWriter)
//       }
//       switch t := err.(type) {
//       default:
//       	w.WriteHeader(status)
//       case HttpError:
//       	error, code := t.HttpError()
//       	t.Error(w, error, code)
//       case HttpError:
//       	t.HttpError(w)
//       }
//   }
//
//   func HandleHttpResponse(w http.ResponseWriter, r *http.Request, resp interface{}) {
//       switch t := resp.(type) {
//       default:
//       	// I don't know that type !
//       case http.Handler:
//       	t.ServeHTTP(w, r)
//       case []byte:
//       	t.Write(t)
//       }
//   }
//
package main // import "github.com/azr/generators/varhandler"

func main() {

}
