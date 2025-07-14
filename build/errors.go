package build

import "errors"

var errCancelled = errors.New("cancelled")

func JoinErrors(err1, err2 error) error {
	if errors.Is(err1, errCancelled) {
		return err2
	}

	if errors.Is(err2, errCancelled) {
		return err1
	}

	return errors.Join(err1, err2)
}
