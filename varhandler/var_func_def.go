package main

import (
	"fmt"
	"go/ast"
	"log"

	_ "golang.org/x/tools/go/gcimporter"
)

//FuncDefinition represents
//the definition of a function
//that's going to be called by the generated code
type FuncDefinition struct {
	Name string // of the function

	//wether or not a status is returned by the handler
	Status bool

	//wether or not a response is returned by the handler
	Response bool

	//params the functions take
	Params []Param
}

type Param struct {
	//Name used inside the function
	Name string

	//the name of the func that will generate our param
	GeneratorName string

	//Defined its a param from another package
	Package string
}

func (fd *FuncDefinition) ParseResults(results *ast.FieldList) bool {
	if len(results.List) == 0 {
		log.Printf("%s should at least return an error", fd.Name)
		return false
	}
	if len(results.List) == 1 {
		return true
	}
	if len(results.List) == 2 {
		v, ok := results.List[0].Type.(*ast.Ident)
		if ok && v.Name == "int" {
			fd.Status = true
		} else {
			fd.Response = true
		}
		return true
	}

	if len(results.List) == 3 {
		fd.Status = true
		fd.Response = true
		return true
	}

	log.Printf("too many results for %s", fd.Name)
	return false
}

func (fd *FuncDefinition) ParseArguments(arguments []*ast.Field) bool {
	generatorNameSuffix := "HTTP"
	for _, argument := range arguments {
		switch v := argument.Type.(type) { // get var type
		case *ast.Ident:
			// plain type like `x X` from `type x struct {}`
			fd.Params = append(fd.Params, Param{
				Name:          v.Name,
				GeneratorName: generatorNameSuffix + v.Name,
			})
		case *ast.StarExpr:
			// arg like `x *X`
			vv, ok := v.X.(*ast.Ident)
			if !ok {
				log.Printf("Found an unary star")
				return false
			}
			fd.Params = append(fd.Params, Param{
				Name:          vv.Name,
				GeneratorName: generatorNameSuffix + vv.Name,
			})
		case *ast.SelectorExpr:
			// arg like `x pkgname.X`
			pkg, ok := v.X.(fmt.Stringer)
			if !ok {
				log.Printf("could not define type of %#v", v)
			}
			fd.Params = append(fd.Params, Param{
				Name:          v.Sel.String(),
				GeneratorName: generatorNameSuffix + v.Sel.String(),
				Package:       pkg.String(),
			})
		default:
			log.Printf("Could not guess var full name, type not expected: %#v", v)
			return false
		}
	}
	return true
}
