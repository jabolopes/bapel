package comp

func interleave[T any](values []T, interleave func(), do func(int, T)) {
	if len(values) == 0 {
		return
	}

	do(0, values[0])
	for i, value := range values[1:] {
		interleave()
		do(i+1, value)
	}
}
