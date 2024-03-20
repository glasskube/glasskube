package semver

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsUpgradable", func() {
	// testCases is a map of installed version to latest version to expected result
	testCases := map[string]map[string]bool{
		"v1.0.0": {
			"v1.0.0":        false,
			"v1.0.0+1":      true,
			"v1.0.0+a":      true,
			"not a version": true,
		},
		"v1.0.0+1": {
			"v1.0.0":        false,
			"v1.0.0+1":      false,
			"v1.0.0+2":      true,
			"v1.0.0+a":      true,
			"not a version": true,
		},
		"v1.0.0+a": {
			"v1.0.0":        false,
			"v1.0.0+1":      true,
			"v1.0.0+a":      false,
			"v1.0.0+b":      true,
			"not a version": true,
		},
		"v1.0.0+b": {
			"v1.0.0":        false,
			"v1.0.0+1":      true,
			"v1.0.0+a":      true,
			"v1.0.0+b":      false,
			"not a version": true,
		},
		"not a version": {
			"v1.0.0":        true,
			"v1.0.0+1":      true,
			"v1.0.0+a":      true,
			"v1.0.0+b":      true,
			"not a version": false,
		},
	}

	for i, v := range testCases {
		// save loop variables locally, this is no longer needed after Go 1.22.0
		installed := i
		When(fmt.Sprintf("installed is %v", installed), func() {
			for l, e := range v {
				latest := l
				expected := e
				When(fmt.Sprintf("latest is %v", latest), func() {
					It(fmt.Sprintf("should return %v", expected), func() {
						Expect(IsUpgradable(installed, latest)).To(Equal(expected))
					})
				})
			}
		})
	}
})
