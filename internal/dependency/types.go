package dependency

import (
	"fmt"
	"strings"
)

type ValidationResultStatus string

const (
	ValidationResultStatusOk         ValidationResultStatus = "OK"
	ValidationResultStatusResolvable ValidationResultStatus = "RESOLVABLE"
	ValidationResultStatusConflict   ValidationResultStatus = "CONFLICT"
)

type PackageWithVersion struct {
	Name    string
	Version string
}

type Requirement struct {
	PackageWithVersion
	ComponentMetadata *ComponentMetadata
	Transitive        bool
}

type ComponentMetadata struct {
	Name, Namespace string
}

type Conflict struct {
	Actual   PackageWithVersion
	Required PackageWithVersion
	Cause    error
}

func (cf Conflict) String() string {
	return fmt.Sprintf("%v (required: %v, actual: %v)", cf.Required.Name, cf.Required.Version, cf.Actual.Version)
}

type Conflicts []Conflict

func (cf Conflicts) String() string {
	s := make([]string, len(cf))
	for i, c := range cf {
		s[i] = c.String()
	}
	return strings.Join(s, ", ")
}

type ValidationResult struct {
	Status       ValidationResultStatus
	Requirements []Requirement
	Conflicts    Conflicts
	Pruned       []Requirement
}
