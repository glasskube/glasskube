package manifestvalues

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateResolvedValues", func() {
	emptyManifest := v1alpha1.PackageManifest{}
	manifestWithSimpleValue := v1alpha1.PackageManifest{
		ValueDefinitions: map[string]v1alpha1.ValueDefinition{
			"test": {Type: v1alpha1.ValueTypeText}},
	}
	manifestWithRequiredValue := v1alpha1.PackageManifest{
		ValueDefinitions: map[string]v1alpha1.ValueDefinition{
			"test": {Type: v1alpha1.ValueTypeText, Constraints: v1alpha1.ValueDefinitionConstraints{Required: true}}},
	}
	five := 5
	ten := 10
	pattern := "a{2,3}b+"
	manifestWithConstraints := v1alpha1.PackageManifest{
		ValueDefinitions: map[string]v1alpha1.ValueDefinition{
			"minmaxstr": {
				Type: v1alpha1.ValueTypeText,
				Constraints: v1alpha1.ValueDefinitionConstraints{
					MinLength: &five,
					MaxLength: &ten,
				},
			},
			"minmax": {
				Type: v1alpha1.ValueTypeNumber,
				Constraints: v1alpha1.ValueDefinitionConstraints{
					Min: &five,
					Max: &ten,
				},
			},
			"pattern": {
				Type: v1alpha1.ValueTypeText,
				Constraints: v1alpha1.ValueDefinitionConstraints{
					Pattern: &pattern,
				},
			},
			"options": {
				Type:    v1alpha1.ValueTypeOptions,
				Options: []string{"foo", "bar"},
			},
			"bool": {Type: v1alpha1.ValueTypeBoolean},
		},
	}
	DescribeTable("Validating values",
		func(manifest v1alpha1.PackageManifest, values map[string]string, valid bool) {
			err := ValidateResolvedValues(manifest, values)
			if valid {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
		},
		Entry("When def is empty and config is empty", emptyManifest, map[string]string{}, true),
		Entry("When value without def", emptyManifest, map[string]string{"test": "test"}, false),
		Entry("When value with matching def", manifestWithSimpleValue, map[string]string{"test": "test"}, true),
		Entry("When required value missing", manifestWithRequiredValue, map[string]string{}, false),
		Entry("When required value present", manifestWithRequiredValue, map[string]string{"test": "test"}, true),
		Entry("When MinLength violated", manifestWithConstraints, map[string]string{"minmaxstr": "aaa"}, false),
		Entry("When MaxLength violated", manifestWithConstraints, map[string]string{"minmaxstr": "aaaaaaaaaaa"}, false),
		Entry("When MinLength, MaxLength not violated",
			manifestWithConstraints, map[string]string{"minmaxstr": "aaaaaaaaaa"}, true),
		Entry("When Min violated", manifestWithConstraints, map[string]string{"minmax": "1"}, false),
		Entry("When Max violated", manifestWithConstraints, map[string]string{"minmax": "11"}, false),
		Entry("When Min, Max not violated", manifestWithConstraints, map[string]string{"minmax": "7"}, true),
		Entry("When Pattern violated", manifestWithConstraints, map[string]string{"pattern": "ab"}, false),
		Entry("When Pattern not violated", manifestWithConstraints, map[string]string{"pattern": "aaab"}, true),
		Entry("When value not in Options", manifestWithConstraints, map[string]string{"options": "test"}, false),
		Entry("When value in Options", manifestWithConstraints, map[string]string{"options": "foo"}, true),
		Entry("When wrong bool format", manifestWithConstraints, map[string]string{"bool": "test"}, false),
		Entry("When correct bool format: true", manifestWithConstraints, map[string]string{"bool": "true"}, true),
		Entry("When correct bool format: false", manifestWithConstraints, map[string]string{"bool": "false"}, true),
		Entry("When correct bool format: 1", manifestWithConstraints, map[string]string{"bool": "1"}, true),
	)
})
