package v1alpha1

import (
	"bytes"

	"github.com/invopop/jsonschema"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type JSON apiextensionsv1.JSON

var nullLiteral = []byte(`null`)

func (JSON) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:                 "object",
		AdditionalProperties: jsonschema.TrueSchema,
	}
}

func (s JSON) MarshalJSON() ([]byte, error) {
	if len(s.Raw) > 0 {
		return s.Raw, nil
	}
	return []byte("null"), nil
}

func (s *JSON) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && !bytes.Equal(data, nullLiteral) {
		s.Raw = data
	}
	return nil
}
