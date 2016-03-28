# varhandler
--
VarHandler generate wrappers for variing http handler funcs

To ease http development process, enabling reusability and remaining http
complient a la go.

Given a pkg :

    	//go:generate varhandler -func F
     package server

     func F(x X, y Y, z *Z, zz z.Z) ([resp interface{},] [status int,] err error) {...}
     // and funcs
     func HTTPX  (r *http.Request) (X, error)   {...}
     func HTTPY  (r *http.Request) (Y, error)   {...}
     func HTTPZ  (r *http.Request) (*Z, error)  {...}
     func z.HTTPZ(r *http.Request) (z.Z, error) {...}

VarHandler will generate an http handler :

    func FVarHandler(w http.ResponseWriter, r *http.Request) {
        var err error
        x, err := HTTPX(r)
        if err != nil {
        	HandleHttpErrorWithDefaultStatus(w, http.StatusBadRequest, err)
        	return
        }
        y, err := HTTPY(r)
        if err != nil {
        	HandleHttpErrorWithDefaultStatus(w, http.StatusBadRequest, err)
        	return
        }
        z, err := HTTPZ(r)
        if err != nil {
        	HandleHttpErrorWithDefaultStatus(w, http.StatusBadRequest, err)
        	return
        }
        zz, err := z.HTTPZ(r)
        if err != nil {
        	HandleHttpErrorWithDefaultStatus(w, http.StatusBadRequest, err)
        	return
        }
        resp, status, err := F(x, y, z, zz)
        if err != nil {
        	HandleHttpErrorWithDefaultStatus(w, http.StatusInternalServerError, err)
        	return
        }
        if status != 0 { // code generated if status is returned
        	w.WriteHeader(status)
        }
        if resp != nil { // code generated if resp object is returned
        	HandleHttpResponse(w, r, resp)
        }
    }

    func HandleHttpErrorWithDefaultStatus(w http.ResponseWriter, status int, err error) {
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
        case SelfHttpError:
        	t.HttpError(w)
        }
    }

    func HandleHttpResponse(w http.ResponseWriter, r *http.Request, resp interface{}) {
        switch t := resp.(type) {
        default:
        	// I don't know that type !
        case http.Handler:
        	t.ServeHTTP(w, r)
        case []byte:
        	w.Write(t)
        }
    }
