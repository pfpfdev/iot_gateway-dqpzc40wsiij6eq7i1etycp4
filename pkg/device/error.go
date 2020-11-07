package device

import (
	"errors"
)

func WrongStateErr(trueState string)error{
	return errors.New("Wrong State("+trueState+")")
}

func WrongFormatErr(trueFormat string)error{
	return errors.New("Wrong Format("+trueFormat+")")
}