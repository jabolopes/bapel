package vm

type bindTable struct {
	ops []Op
}

func newBindTable() bindTable {
	return bindTable{NewOpTable().ops}
}
