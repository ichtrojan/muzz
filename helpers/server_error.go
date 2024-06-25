package helpers

import "errors"

func ServerError() error {
	return errors.New("something went wrong")
}
