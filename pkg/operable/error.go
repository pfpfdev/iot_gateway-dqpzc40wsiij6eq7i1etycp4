package operable

import (
	"errors"
)

func UndefinedTypeErr() error{
	return errors.New("Undefind Operation Type")
} 
func UndefinedOperationErr() error{
	return errors.New("Undefind Operation")
} 
func InvalidArgumentErr()error{
	return errors.New("Invalid Operation Argument")
}