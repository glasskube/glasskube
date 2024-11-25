package components

import (
	"fmt"
)

type DatalistInput struct {
	Id      string
	Options []string
}

func ForDatalist(valueName string, postfix string, options []string) DatalistInput {
	if postfix != "" {
		postfix = "-" + postfix
	}
	return DatalistInput{
		Id:      fmt.Sprintf("%s%s", valueName, postfix),
		Options: options,
	}
}
