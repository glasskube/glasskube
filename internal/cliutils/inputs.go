package cliutils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/glasskube/glasskube/internal/config"
)

// InteractivityEnabledOrFail checks whether config.NonInteractive is set and immediately aborts if it is not.
func InteractivityEnabledOrFail() {
	if config.NonInteractive {
		fmt.Fprintln(os.Stderr, "\n❌ Interactivity was requested in non-interactive mode")
		ExitWithError()
	}
}

func YesNoPrompt(label string, defaultChoice bool) bool {
	choices := "Y/n"
	if !defaultChoice {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string
	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", strings.TrimSpace(label), choices)
		InteractivityEnabledOrFail()
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return defaultChoice
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func GetInputStr(label string) (input string) {
	fmt.Fprintf(os.Stderr, "%v> ", label)
	InteractivityEnabledOrFail()
	fmt.Scanln(&input)
	return
}

func GetOption(label string, options []string) (string, error) {
	return GetOptionWithDefault(label, options, nil)
}

func GetOptionWithDefault(label string, options []string, defaultOption *string) (string, error) {
	printOptions(options, defaultOption)
	iStr := GetInputStr(label)
	if len(iStr) == 0 && defaultOption != nil {
		return *defaultOption, nil
	} else if i, err := strconv.Atoi(iStr); err != nil {
		return "", err
	} else {
		i-- // index is entered with offset 1
		if 0 <= i && i < len(options) {
			return options[i], nil
		} else {
			return "", fmt.Errorf("%v is not a valid option", i+1)
		}
	}
}

func printOptions(options []string, defaultOption *string) {
	msg := "Enter the number of one of the following"
	if defaultOption != nil {
		msg += fmt.Sprintf(" (default: %v)", *defaultOption)
	}
	fmt.Fprintln(os.Stderr, msg+":")
	for i, opt := range options {
		fmt.Fprintf(os.Stderr, "%v) %v\n", i+1, opt)
	}
}
