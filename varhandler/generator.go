// Generate boilerplate http input and ouput for golang
//
// To ease http development process,
// enabling reusability
// and remaining http complient a la go.
//
// In short
//
//
// Golang HTTP dev at its essense will be
//  1/ get parameters from http request & validate them ( n times )
//  2/ if something is wrong do something about it      ( n times )
//  3/ do the job given set of parameters               (from steps 1/)
//  4/ if something went wrong do something about it
//  5/ http respond
//
// Now given 3 (a variing function that does a job),
// and 1 (n instantiator/validators)
//
// varhandler will generate steps 1, 2, 4 and 5
//
// Variing types of response
//
// The function must return an error.
// Otherwise return can take multiple forms:
//
//  func F(x X, y Y) (err error)                                   // default return status when no error is 200 - OK
//                                                                 // if the error is unknown the status will be
//                                                                 // InternalServerError - see HandleHttpErrorWithDefaultStatus
//                                                                 // for error types
//
//  func F(x X, y Y) (status int, err error)                       // sets returned status if error is nil
//
//  func F(x X, y Y) (response interface{}, err error)             // specific response See HandleHttpResponse code
//
//  func F(x X, y Y) (response interface{}, status int, err error) // sets status and does Response Handling if no error is set
//
// Variing parameters
//
// The functions takes one or more arguments.
// Those arguments need to have http instantiators
//  HTTPX(r *http.Request) (x X, err error)
//
// Error handling
//
// If an instantiation error occurs:
//  HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err) // will be called
//
// If the wrapped func returns an error
//  HandleHttpErrorWithDefaultStatus(w, r, http.StatusInternalServerError, err) // will be called
//
// Response handling
//
// check HandleHttpResponse's code
//
// Example
//
// Old way :
//
//  myHttpHandler(w http.ResponseWriter, r *http.Request) {
//      var x X
//      err := json.NewDecoder(r.Body).Decode(x)
//      if err != nil {
//         //do something about it
//         return
//      }
//      err := validate(x)
//      if err != nil {
//         //do something about it
//         return
//      }
//      // repeat for y and z
//      // get zz from zz pkg
//      response, status, err := F(x, y, z, zz)
//      if err != nil {
//          // return http.StatusInternalServerError
//      }
//      w.WriteHeader(status)
//      json.NewEncoder(w).Encode(response)
//  }
//
// Now The only interesting part should be the call to F.
//
// New way:
//  //go:generate varhandler -func F
//  package server
//
//  import "github.com/azr/generators/varhandler/examples/z"
//
//  func F(x X, y Y, z *Z, zz z.Z) (resp interface{}, status int, err error) {...}
//
//  // http instantiators :
//
//  func HTTPX  (r *http.Request) (x X, err error)   { err = json.NewDecoder(r.Body).Decode(&x) }
//  func HTTPY  (r *http.Request) (Y, error)         {...}
//  func HTTPZ  (r *http.Request) (*Z, error)        {...}
//  func z.HTTPZ(r *http.Request) (z.Z, error)       {...}
//
//will be generated:
//
// * calls to HTTPX, HTTPY, HTTPX and z.HTTPZ and their error check
//
// * call to F(x, y, z, zz)
//
// * http response code given F's return arguments
//
//
// In that case generated code looks like :
//
//   func FHandler(w http.ResponseWriter, r *http.Request) {
//       var err error
//       x, err := HTTPX(r)
//       if err != nil {
//          HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
//          return
//       }
//       y, err := HTTPY(r)
//       if err != nil {
//          HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
//          return
//       }
//       z, err := HTTPZ(r)
//       if err != nil {
//          HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
//          return
//       }
//       zz, err := z.HTTPZ(r)
//       if err != nil {
//          HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
//          return
//       }
//       resp, status, err := F(x, y, z, zz)
//       if err != nil {
//          HandleHttpErrorWithDefaultStatus(w, r, http.StatusInternalServerError, err)
//          return
//       }
//       if status != 0 { // code generated if status is returned by F
//          w.WriteHeader(status)
//       }
//       if resp != nil { // code generated if resp object is returned by F
//          HandleHttpResponse(w, r, resp)
//       }
//   }
//
//   //Helper funcs
//
//   func HandleHttpErrorWithDefaultStatus(w http.ResponseWriter, r *http.Request, status int, err error) {
//       type HttpError interface {
//          HttpError() (error string, code int)
//       }
//       type SelfHttpError interface {
//          HttpError(w http.ResponseWriter)
//       }
//       switch t := err.(type) {
//       default:
//          w.WriteHeader(status)
//       case HttpError:
//          err, code := t.HttpError()
//          http.Error(w, err, code)
//       case http.Handler:
//          t.ServeHTTP(w, r)
//       case SelfHttpError:
//          t.HttpError(w)
//       }
//   }
//
//   func HandleHttpResponse(w http.ResponseWriter, r *http.Request, resp interface{}) {
//       type Byter interface {
//          Bytes() []byte
//       }
//       type Stringer interface {
//          String() string
//       }
//       switch t := resp.(type) {
//       default:
//          // I don't know that type !
//       case http.Handler:
//          t.ServeHTTP(w, r) // resp knows how to handle itself
//       case Byter:
//          w.Write(t.Bytes())
//       case Stringer:
//          w.Write([]byte(t.String()))
//       case []byte:
//          w.Write(t)
//       }
//   }
//
package main // import "github.com/azr/generators/varhandler"

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/types"

	_ "golang.org/x/tools/go/gcimporter"

	"github.com/azr/generators/utils"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tvarhandler [flags] -func F [directory]\n")
	fmt.Fprintf(os.Stderr, "\tvarhandler [flags] -func F files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttp://godoc.org/github.com/azr/generators/varhandler\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	{ //setup logs
		log.SetFlags(0)
		log.SetPrefix("handler: ")
	}

	var funcNames, output string
	{ // init
		flag.StringVar(&funcNames, "func", "", "comma-separated list of func names; must be set")
		flag.StringVar(&output, "output", "", "output file name; default srcdir/generated_varhandlers.go")
		flag.Usage = Usage
		flag.Parse()
	}
	if len(funcNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	funcs := strings.Split(funcNames, ",")

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var (
		dir string
		g   Generator
	)
	if len(args) == 1 && utils.IsDirectory(args[0]) {
		dir = args[0]
		g.parsePackageDir(args[0])
	} else {
		dir = filepath.Dir(args[0])
		g.parsePackageFiles(args)
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"handler %s\"; DO NOT EDIT\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s\n", g.pkg.name)
	g.Printf("\n")
	g.Printf("import \"net/http\"\n") // Used by all methods.

	var definitions []FuncDefinition

	for _, funcName := range funcs {
		// generate import for func if any
		// and generate definition of func for latter call
		definitions = append(definitions, g.generateImportPaths(funcName))
	}
	for _, definition := range definitions {
		if definition.Name != "" { // func was found
			log.Printf("Defining: %s", definition.Name)
			g.writeFuncDef(definition)
		}
	}

	g.Printf(utilFuncs)

	// Format the output.
	src := g.format()

	// Write to file.
	outputName := output
	if outputName == "" {
		outputName = filepath.Join(dir, "generated_varhandlers.go")
	}
	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.

	// These fields are reset for each func being generated.
	funcDefinition FuncDefinition
	found          bool
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	pkgs     map[string]*types.Package
	files    []*File
	typesPkg *types.Package
}

// parsePackageDir parses the package residing in the directory.
func (g *Generator) parsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var names []string
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	// TODO: Need to think about constants in test files. Maybe write type_string_test.go
	// in a separate pass? For later.
	// names = append(names, pkg.TestGoFiles...) // These are also in the "foo" package.
	names = append(names, pkg.SFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(directory, names, nil)
}

// parsePackageFiles parses the package occupying the named files.
func (g *Generator) parsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
}

// prefixDirectory places the directory name on the beginning of each name in the list.
func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

// parsePackage analyzes the single package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing. parsePackage exits if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, text, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}

	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)
}

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{
		FakeImportC: true,
		Packages:    make(map[string]*types.Package),
	}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
	pkg.pkgs = config.Packages
}

// generateImportPaths parses the funcs that are going to be called
// and generates import paths if any generator is in another pkg
func (g *Generator) generateImportPaths(funcName string) FuncDefinition {
	found := false
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.found = false
		file.funcDefinition = FuncDefinition{
			Name: funcName,
		}
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			if file.found {
				for _, param := range file.funcDefinition.Params {
					if param.Package != "" {
						for path, pkg := range g.pkg.pkgs {
							if pkg.Name() == param.Package {
								g.Printf("import %s \"%s\"\n", param.Package, path)
							}
						}
					}
				}

				found = true
				return file.funcDefinition
			}
		}
	}

	if !found {
		fmt.Printf("Func not found: %s", funcName)
	}
	return FuncDefinition{}
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.FuncDecl)
	if !ok {
		// We only care about func declarations.
		return true
	}
	if decl.Name.Name == f.funcDefinition.Name {
		if len(decl.Type.Params.List) == 0 {
			log.Printf("%s should take at least one parameter, found %d instead", f.funcDefinition.Name, len(decl.Type.Params.List))
			return false
		}
		ok := f.funcDefinition.ParseResults(decl.Type.Results)
		if ok {
			ok = f.funcDefinition.ParseArguments(decl.Type.Params.List)
		}

		f.found = ok
	}
	return false
}

// writeFuncDef generates an handler func
func (g *Generator) writeFuncDef(fd FuncDefinition) {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t := template.Must(template.New("varhandler").Funcs(funcMap).Parse(handlerWrap))

	err := t.Execute(&g.buf, fd)
	checkError(err)
}

const handlerWrap = `
func {{.Name}}Handler(w http.ResponseWriter, r *http.Request) {
    var err error
{{range $i, $param := .Params}}
    param{{$i}}, err := {{if ne $param.Package ""}}{{$param.Package}}.{{end}}{{$param.GeneratorName}}(r)
    if err != nil {
        HandleHttpErrorWithDefaultStatus(w, r, http.StatusBadRequest, err)
        return
    }
{{end}}
{{if .Response}}
    var resp interface{}
{{end}}
{{if .Status}}
    var status int
{{end}}
    {{if .Response}}resp, {{end}}{{if .Status}}status, {{end}}err = {{.Name}}({{range $i, $param := .Params}} {{if gt $i 0}},{{end}} param{{$i}}{{end}})
    if err != nil {
        HandleHttpErrorWithDefaultStatus(w, r, http.StatusInternalServerError, err)
        return
    }
{{if .Status}}
    if status != 0 {
        w.WriteHeader(status)
    }
{{end}}
{{if .Response}}
    if resp != nil {
        HandleHttpResponse(w, r, resp)
    }
{{end}}
}
`

const utilFuncs = `
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
`

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
