package shift

import "strings"

func TrimPrefix(arg *string, token string, err error) error {
	if !strings.HasPrefix(*arg, token) {
		return err
	}
	*arg = strings.TrimPrefix(*arg, token)
	return nil
}

func TrimSuffix(arg *string, token string, err error) error {
	if !strings.HasSuffix(*arg, token) {
		return err
	}
	*arg = strings.TrimSuffix(*arg, token)
	return nil
}
