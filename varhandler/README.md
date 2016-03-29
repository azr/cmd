# varhandler
--
Generate boilerplate http input and ouput for golang

To ease http development process, enabling reusability and remaining http
complient a la go.


In short

Golang HTTP dev at its essense will be

    1/ get parameters from http request & validate them ( n times )
    2/ if something is wrong do something about it      ( n times )
    3/ do the job given set of parameters               (from steps 1/)
    4/ if something went wrong do something about it
    5/ http respond

Now given 3 (a variing function that does a job), and 1 (n
instantiator/validators)

varhandler will generate steps 1, 2, 4 and 5


Variing types of response

The function must return an error. Otherwise return can take multiple forms:

    func F(x X, y Y) (err error)                                   // default return status when no error is 200 - OK
                                                                   // if the error is unknown the status will be
                                                                   // InternalServerError - see HandleHttpErrorWithDefaultStatus
                                                                   // for error types

    func F(x X, y Y) (status int, err error)                       // sets returned status if error is nil

    func F(x X, y Y) (response interface{}, err error)             // specific response See HandleHttpResponse code

    func F(x X, y Y) (response interface{}, status int, err error) // sets status and does Response Handling if no error is set


Variing parameters

The functions takes one or more arguments. Those arguments need to have http
### instantiators

    HTTPX(r *http.Request) (x X, err error)


Error handling

If an instantiation error occurs:

    HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err) // will be called

If the wrapped func returns an error

    HandleHttpErrorWithDefaultStatus(w, r, http.StatusInternalServerError, err) // will be called


Response handling

check HandleHttpResponse's code


### Example

Old way :

    myHttpHandler(w http.ResponseWriter, r *http.Request) {
        var x X
        err := json.NewDecoder(r.Body).Decode(x)
        if err != nil {
           //do something about it
           return
        }
        err := validate(x)
        if err != nil {
           //do something about it
           return
        }
        // repeat for y and z
        // get zz from zz pkg
        response, status, err := F(x, y, z, zz)
        if err != nil {
            // return http.StatusInternalServerError
        }
        w.WriteHeader(status)
        json.NewEncoder(w).Encode(response)
    }

Now The only interesting part should be the call to F.

New way:

    //go:generate varhandler -func F
    package server

    import "github.com/azr/generators/varhandler/examples/z"

    func F(x X, y Y, z *Z, zz z.Z) (resp interface{}, status int, err error) {...}

    // http instantiators :

    func HTTPX  (r *http.Request) (x X, err error)   { err = json.NewDecoder(r.Body).Decode(&x) }
    func HTTPY  (r *http.Request) (Y, error)         {...}
    func HTTPZ  (r *http.Request) (*Z, error)        {...}
    func z.HTTPZ(r *http.Request) (z.Z, error)       {...}

will be generated:

* calls to HTTPX, HTTPY, HTTPX and z.HTTPZ and their error check

* call to F(x, y, z, zz)

* http response code given F's return arguments

In that case generated code looks like :

    func FHandler(w http.ResponseWriter, r *http.Request) {
        var err error
        x, err := HTTPX(r)
        if err != nil {
           HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
           return
        }
        y, err := HTTPY(r)
        if err != nil {
           HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
           return
        }
        z, err := HTTPZ(r)
        if err != nil {
           HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
           return
        }
        zz, err := z.HTTPZ(r)
        if err != nil {
           HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
           return
        }
        resp, status, err := F(x, y, z, zz)
        if err != nil {
           HandleHttpErrorWithDefaultStatus(w, r, http.StatusInternalServerError, err)
           return
        }
        if status != 0 { // code generated if status is returned by F
           w.WriteHeader(status)
        }
        if resp != nil { // code generated if resp object is returned by F
           HandleHttpResponse(w, r, resp)
        }
    }

    //Helper funcs

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
