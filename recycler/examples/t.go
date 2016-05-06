package examples

//go:generate recycler -type=T -output t_pool.go
//go:generate recycler -type=T -template freelist.gotpl -output t_freelist.go
type T struct {
}
