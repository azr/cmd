#Generators

Generate go files with common and repetitive usages.

## *[varhandler](/varhandler)*

Generate wrappers for variing http handler funcs.

Given
```
func F(x X, y Y, z *Z, zz z.Z) (err error) {...}
//or
func F(x X, y Y, z *Z, zz z.Z) (status int, err error) {...}
//or
func F(x X, y Y, z *Z, zz z.Z) (resp interface{}, status int, err error) {...}
//or
func F(x X, y Y, z *Z, zz z.Z) (resp http.Handler, status int, err error) {...}
//or
func F(x X, y Y, z *Z, zz z.Z) (resp []byte, status int, err error) {...}
```
Generate the code that will:
* instantiate x, y, z and zz
* call F(x, y, z, zz)
* check return statuses



## *[pooler](/pooler)*

Generate *typed* [sync.Pool](https://golang.org/pkg/sync/#Pool) wrappers

## *[handler](/handler)*
Generate typed http handlers

From

[http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)

to
```
func ([context.Context,] YourStruct) (resp interface{}, status int)
```
Generator
