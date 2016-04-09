package main

import "github.com/azr/generators/varhandler/examples/z"

//pkg z defines an HTTPZ function
//go:generate varhandler -func Import
func Import(zz z.Z) error {
	return nil
}
