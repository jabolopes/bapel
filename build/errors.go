package build

import (
	"context"
	"errors"
)

var errCancelled = context.Canceled

func JoinErrors(err1, err2 error) error {
	if errors.Is(err1, errCancelled) && err2 != nil {
		return err2
	}

	if errors.Is(err2, errCancelled) && err1 != nil {
		return err1
	}

	return errors.Join(err1, err2)
}
