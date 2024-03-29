package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	corev1 "k8s.io/api/core/v1"
)

// +kubebuilder:validation:Enum=boolean;text;number;options
type ValueType string

const (
	ValueTypeBoolean ValueType = "boolean"
	ValueTypeText    ValueType = "text"
	ValueTypeNumber  ValueType = "number"
	ValueTypeOptions ValueType = "options"
)

func (ref *ValueType) parseString(data string) error {
	switch data {
	case string(ValueTypeBoolean):
		*ref = ValueTypeBoolean
	case string(ValueTypeText):
		*ref = ValueTypeText
	case string(ValueTypeNumber):
		*ref = ValueTypeNumber
	case string(ValueTypeOptions):
		*ref = ValueTypeOptions
	default:
		return fmt.Errorf("invalid ValueType: %v", data)
	}
	return nil
}
func (ValueType) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{ValueTypeBoolean, ValueTypeText, ValueTypeNumber, ValueTypeOptions},
	}
}

func (ref *ValueType) UnmarshalJson(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return ref.parseString(s)
}

type ValueDefinitionMetadata struct {
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	Hints       []string `json:"hints,omitempty"`
}

type ValueDefinitionConstraints struct {
	Required  bool `json:"required,omitempty"`
	Min       int  `json:"min,omitempty"`
	Max       int  `json:"max,omitempty"`
	MinLength int  `json:"minLength,omitempty"`
	MaxLength int  `json:"maxLength,omitempty"`
}

type PartialJsonPatch struct {
	Op   string `json:"op" jsonschema:"required"`
	Path string `json:"path" jsonschema:"required"`
}

// +kubebuilder:validation:XValidation:message="ValueDefinitionTarget must have either resource or chartName but not both",rule="has(self.resource) != has(self.chartName)"
type ValueDefinitionTarget struct {
	Resource      *corev1.TypedObjectReference `json:"resource,omitempty" jsonschema:"oneof_required=WithResource"`
	ChartName     *string                      `json:"chartName,omitempty" jsonschema:"oneof_required=WithChartName"`
	Patch         PartialJsonPatch             `json:"patch" jsonschema:"required"`
	ValueTemplate string                       `json:"valueTemplate,omitempty"`
}

type ValueDefinition struct {
	Type         ValueType                  `json:"type" jsonschema:"required"`
	Metadata     ValueDefinitionMetadata    `json:"metadata,omitempty"`
	DefaultValue string                     `json:"defaultValue,omitempty"`
	Options      []string                   `json:"options,omitempty"`
	Constraints  ValueDefinitionConstraints `json:"constraints,omitempty"`
	Targets      []ValueDefinitionTarget    `json:"targets" jsonschema:"required"`
}
