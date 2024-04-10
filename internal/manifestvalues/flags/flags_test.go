package flags

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseValues", func() {
	foo := "foo"
	fooMap := map[string]v1alpha1.ValueConfiguration{"foo": {Value: &foo}}
	var opts *ValuesOptions
	BeforeEach(func() { opts = &ValuesOptions{} })
	DescribeTable("should parse",
		func(keep bool, values []string, oldValues map[string]v1alpha1.ValueConfiguration,
			expectError bool, expectedResult map[string]v1alpha1.ValueConfiguration) {
			opts.KeepOldValues = keep
			opts.Values = values
			newValues, err := opts.ParseValues(oldValues)
			if expectError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
			if expectedResult == nil {
				Expect(newValues).To(BeNil())
			} else {
				Expect(newValues).To(Equal(expectedResult))
			}
		},

		Entry("when KeepOldValues is true", true, nil, fooMap, false, fooMap),

		Entry("when there is an invalid value", false, []string{"foo"}, nil, true, nil),

		Entry("when there is a literal value", false, []string{"foo=foo"}, nil, false, fooMap),

		Entry("when there is a valid ConfigMapRef", false, []string{"foo=$ConfigMapRef$foo,bar,data"}, nil, false,
			map[string]v1alpha1.ValueConfiguration{
				"foo": {ValueFrom: &v1alpha1.ValueReference{ConfigMapRef: &v1alpha1.ObjectKeyValueSource{
					Namespace: "foo",
					Name:      "bar",
					Key:       "data",
				}}},
			}),

		Entry("when there is an invalid ConfigMapRef (too few args)", false, []string{"foo=$ConfigMapRef$foo,bar"},
			nil, true, nil),

		Entry("when there is a valid SecretRef", false, []string{"foo=$SecretRef$foo,bar,data"}, nil, false,
			map[string]v1alpha1.ValueConfiguration{
				"foo": {ValueFrom: &v1alpha1.ValueReference{SecretRef: &v1alpha1.ObjectKeyValueSource{
					Namespace: "foo",
					Name:      "bar",
					Key:       "data",
				}}},
			}),

		Entry("when there is an invalid SecretRef (too few args)", false, []string{"foo=$SecretRef$foo,bar"}, nil, true, nil),

		Entry("when there is a valid PackageRef", false, []string{"foo=$PackageRef$foo,data"}, nil, false,
			map[string]v1alpha1.ValueConfiguration{
				"foo": {ValueFrom: &v1alpha1.ValueReference{PackageRef: &v1alpha1.PackageValueSource{
					Name:  "foo",
					Value: "data",
				}}},
			}),

		Entry("when there is an invalid PackageRef (too few args)", false, []string{"foo=$PackageRef$foo"}, nil, true, nil),
	)
})
