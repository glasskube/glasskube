package plain

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getActualManifestUrl", func() {
	It("should handle relative url", func() {
		result, err := getActualManifestUrl(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"./manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("http://localhost/packages/foo/manifest.yaml"))
	})
	It("should handle plain file name", func() {
		result, err := getActualManifestUrl(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("http://localhost/packages/foo/manifest.yaml"))
	})
	It("should handle relative url with \"..\"", func() {
		result, err := getActualManifestUrl(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"../manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("http://localhost/packages/manifest.yaml"))
	})
	It("should handle absolute path", func() {
		result, err := getActualManifestUrl(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"/manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("http://localhost/manifest.yaml"))
	})
	It("should handle real url", func() {
		result, err := getActualManifestUrl(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"https://github.com/glasskube/glasskube/manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("https://github.com/glasskube/glasskube/manifest.yaml"))
	})
})
