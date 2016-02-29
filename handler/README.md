# handler
--
Handler builds typed golang http handlers.

Given a func F :

    func F(x X) (status int, resp interface{}) )

and a format pkg like encoding/json.

handler will create an http handler :

    func FHandlerFORMAT(w http.ResponseWriter, r *http.Request)
    // decode
    // call F()

The file is created in the same package and directory as the package that
defines F.

ex:

    //go:generate handler -func=PutJob -format=json
    package jober

    type job struct { A string }

    func PutJob(j job) (int, interface{}) {
      return nil, 200
    }

### running

    go generate pkg.go/foo/jober

will create generated_handlers.go:

    import "encoding/json"

    func PutJobHandlerJSON(w http.ResponseWriter, r *http.Request) {
        x := job{}
        err := json.NewDecoder(r.Body).Decode(&x)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        s, resp := PutJob(x)
        w.WriteHeader(s)
        json.NewEncoder(w).Encode(resp)
    }

so now you can just worry about what PutJob does.

pkg existence will be checked. The pkg needs to have funcs :

    func NewDecoder(r io.Reader) *Decoder
    func NewEncoder(r io.Reader) *Encoder

and types

    type Encoder interface {
        Encode(v interface{}) error
    }
    type Decoder interface {
        Decode(v interface{}) error
    }

Typically this process would be run using go generate, by writing:

    //go:generate handler -encoding encoding/json -func PutJob

at the beginning of your .go file

The -encoding and the -func flags accepts a comma-separated list of strings. So
you can have n handler working in m encoding

Name of the created file can be overridden with the -output flag.

Support of contexts is comming soon.
