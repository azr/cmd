package main

//go:generate varhandler -func Simple
func Simple(x X, y Y, z *Z) error {
	return nil
}
