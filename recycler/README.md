# Recycler
recycler is a tool to automate the creation of typed memory recyclers in golang

[![GoDoc](https://godoc.org/github.com/azr/generators/recycler?status.png)](https://godoc.org/github.com/azr/generators/recycler)

You can generate multiple type of recycler, they all have the same signature :

* To generate pool backed recycler: 
```
//go:generate recycler -type=<T> -output <file.go>
```

* To generate freelist backed recycler: 
```
//go:generate recycler -type=<T> -size <buffer_size> -template freelist.gotpl
```

