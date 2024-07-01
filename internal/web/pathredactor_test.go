package web

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path Redactor", func() {
	DescribeTable("Server Paths",
		func(url string, result string) {
			Expect(packagesPathRedactor(url)).To(Equal(result))
		},
		Entry("Empty Path", "/", "/"),
		Entry("Packages Path", "/packages", "/packages"),
		Entry("Clusterpackages Path", "/clusterpackages", "/clusterpackages"),
		Entry("Settings Path", "/settings", "/settings"),
		Entry("Invalid Package Path", "/packages/manifest", "/packages/manifest"),
		Entry("Installed Package Path", "/packages/manifest/namespace/name", "/packages/manifest/x/x"),
		Entry("Installed Package Path with query params", "/packages/manifest/namespace/name?someVar=someValue&x=y",
			"/packages/manifest/x/x?someVar=someValue&x=y"),
		Entry("Installed Package Discussion Path with query params", "/packages/manifest/namespace/name/discussion?someVar=someValue&x=y",
			"/packages/manifest/x/x/discussion?someVar=someValue&x=y"),
		Entry("Installed Package Configure Path", "/packages/manifest/namespace/name/configure", "/packages/manifest/x/x/configure"),
		Entry("Uninstalled Package Configure Path", "/packages/manifest/-/-/configure", "/packages/manifest/-/-/configure"),
		Entry("Invalid Package Path", "/packages/manifest/whatever", "/packages/manifest/whatever"),
		Entry("Invalid Package Path with query params", "/packages/manifest/whatever?someVar=someValue&x=y",
			"/packages/manifest/whatever?someVar=someValue&x=y"),
	)
})
