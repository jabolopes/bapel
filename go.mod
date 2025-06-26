module github.com/jabolopes/bapel

go 1.23

toolchain go1.23.4

require (
	github.com/emirpasic/gods/v2 v2.0.0-alpha
	github.com/golang/glog v1.2.4
	github.com/google/go-cmp v0.6.0
	github.com/jabolopes/go-lalr1 v0.0.0-00010101000000-000000000000
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842
)

require github.com/deckarep/golang-set/v2 v2.6.0 // indirect

replace github.com/jabolopes/go-lalr1 => /home/jose/Projects/go/src/github.com/jabolopes/go-lalr1
