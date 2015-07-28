package main

import (
	"errors"

	//"github.com/ActiveState/tail"
)

func TailValidateConfig(config *ConfStruct) error {
	if config.TailPath == "" {
		return errors.New("input mode 'tail' requires VEILLE_TAIL_PATH environment variable")
	}
	return nil
}
