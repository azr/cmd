#Generators

Generate go files with common and repetitive usages.

## *[varhandler](/varhandler)*

Generate wrappers for variing http handler funcs.

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
