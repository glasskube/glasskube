package components

import (
	"fmt"
	"github.com/glasskube/glasskube/internal/web/types"
)

type DatalistInput struct {
	types.TemplateContextHolder
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
