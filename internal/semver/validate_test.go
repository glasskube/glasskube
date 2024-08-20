package semver

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidateConstraint", func() {
	DescribeTable("Checking version ranges",
		func(versionRange string, version string, valid bool) {
			err := ValidateConstraint(version, versionRange)
			if valid {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError((error)(&ConstraintValidationError{})))
			}
		},
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.2.0", true),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.2.0+0", true),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.2.4", true),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.2.9", true),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.1.9", false),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.2.0-rc.1", false),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.3.0-alpha.1", false),
		Entry("When minor version is pinned", ">=1.2.0 <1.3.0", "1.3.0", false),

		Entry("When minor version is pinned with ~", "~1.2.0", "1.2.0", true),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.2.0+0", true),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.2.4", true),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.2.9", true),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.1.9", false),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.2.0-rc.1", false),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.3.0-alpha.1", false),
		Entry("When minor version is pinned with ~", "~1.2.0", "1.3.0", false),

		Entry("When minor version is pinned with ~", "~1.2", "1.2.0", true),
		Entry("When minor version is pinned with ~", "~1.2", "1.2.0+0", true),
		Entry("When minor version is pinned with ~", "~1.2", "1.2.4", true),
		Entry("When minor version is pinned with ~", "~1.2", "1.2.9", true),
		Entry("When minor version is pinned with ~", "~1.2", "1.1.9", false),
		Entry("When minor version is pinned with ~", "~1.2", "1.2.0-rc.1", false),
		Entry("When minor version is pinned with ~", "~1.2", "1.3.0-alpha.1", false),
		Entry("When minor version is pinned with ~", "~1.2", "1.3.0", false),

		Entry("When minor version is pinned with .x", "1.2.x", "1.2.0", true),
		Entry("When minor version is pinned with .x", "1.2.x", "1.2.0+0", true),
		Entry("When minor version is pinned with .x", "1.2.x", "1.2.4", true),
		Entry("When minor version is pinned with .x", "1.2.x", "1.2.9", true),
		Entry("When minor version is pinned with .x", "1.2.x", "1.1.9", false),
		Entry("When minor version is pinned with .x", "1.2.x", "1.2.0-rc.1", false),
		Entry("When minor version is pinned with .x", "1.2.x", "1.3.0-alpha.1", false),
		Entry("When minor version is pinned with .x", "1.2.x", "1.3.0", false),

		Entry("When major version is pinned with ^", "^1", "1.0.0", true),
		Entry("When major version is pinned with ^", "^1", "1.3.5", true),
		Entry("When major version is pinned with ^", "^1", "1.7.9+1", true),
		Entry("When major version is pinned with ^", "^1", "1.0.0-rc.1", false),
		Entry("When major version is pinned with ^", "^1", "2.0.0-rc.1", false),
		Entry("When major version is pinned with ^", "^1", "2.0.0", false),

		Entry("When major version is pinned with .x", "1.x", "1.0.0", true),
		Entry("When major version is pinned with .x", "1.x", "1.3.5", true),
		Entry("When major version is pinned with .x", "1.x", "1.7.9+1", true),
		Entry("When major version is pinned with .x", "1.x", "1.0.0-rc.1", false),
		Entry("When major version is pinned with .x", "1.x", "2.0.0-rc.1", false),
		Entry("When major version is pinned with .x", "1.x", "2.0.0", false),

		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "1.0.0-alpha.1", true),
		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "1.0.0-beta.0", true),
		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "1.0.0-rc.0", true),
		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "1.0.0+0", true),
		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "3.0.0", true),
		Entry("When a prerelease is the minimum", ">= 1.0.0-alpha.1", "1.0.0-alpha.0", false),

		Entry("When a minimum contains a build number", ">= 1.0.0+3", "1.0.0+3", true),
		Entry("When a minimum contains a build number", ">= 1.0.0+3", "1.0.0+2", true), // TODO to be fixed with #405
	)
})
