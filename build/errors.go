package build

import (
	"context"
	"errors"
)

var errCancelled = context.Canceled

func JoinErrors(err1, err2 error) error {
	if errors.Is(err1, errCancelled) {
		return err2
	}

	if errors.Is(err2, errCancelled) {
		return err1
	}

	return errors.Join(err1, err2)
}
