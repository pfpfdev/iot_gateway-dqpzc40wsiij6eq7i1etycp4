package device

import (
	"regexp"
)

type Operation struct{
	Cmd string
	Type string
}

var definedType = map[string]*regexp.Regexp{
	"OnOff":regexp.MustCompile("^[On|Off]$"),
	"Hundred":regexp.MustCompile("^[0-9]{1,2}$"),
}