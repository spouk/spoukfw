package spoukfw

import (
	"fmt"
	"errors"
)

func makeErrorMessage(def, message string) error {
	return errors.New(fmt.Sprintf(def, message))
}
//externExport
func MakeErrorMessage(def, message string) error {
	return errors.New(fmt.Sprintf(def, message))
}

