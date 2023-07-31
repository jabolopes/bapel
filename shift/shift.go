package shift

func Shift(args []string, err error) (string, []string, error) {
	if len(args) == 0 {
		return "", nil, err
	}
	return args[0], args[1:], nil
}

func ShiftIf(args []string, token string, err error) ([]string, error) {
	if len(args) == 0 || args[0] != token {
		return nil, err
	}
	return args[1:], nil
}

func ShiftIfEnd(args []string, token string, err error) ([]string, error) {
	if len(args) == 0 || args[len(args)-1] != token {
		return nil, err
	}
	return args[:len(args)-1], nil
}
