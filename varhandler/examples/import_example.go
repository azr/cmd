//go:generate varhandler -func Simple,Import,Status,Response,ResponseStatus
package main

import "github.com/azr/generators/varhandler/examples/z"

//z.Z defines the HTTPZ function
func Import(x X, y Y, z *Z, zz z.Z) error {
	return nil
}
