package plain

import (
	"github.com/glasskube/glasskube/api/v1alpha1"
	"github.com/glasskube/glasskube/internal/repo/client/auth"
	"github.com/glasskube/glasskube/internal/repo/client/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var adapter = &Adapter{}

var _ = Describe("getActualManifestUrl", func() {
	BeforeEach(func() {
		adapter.repo = fake.EmptyClientset()
	})

	It("should handle relative url", func() {
		result, err := adapter.newManifestRequest(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"./manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.URL.String()).To(Equal("http://localhost/packages/foo/manifest.yaml"))
	})
	It("should handle plain file name", func() {
		result, err := adapter.newManifestRequest(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.URL.String()).To(Equal("http://localhost/packages/foo/manifest.yaml"))
	})
	It("should handle relative url with \"..\"", func() {
		result, err := adapter.newManifestRequest(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"../manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.URL.String()).To(Equal("http://localhost/packages/manifest.yaml"))
	})
	It("should handle absolute path", func() {
		result, err := adapter.newManifestRequest(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"/manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.URL.String()).To(Equal("http://localhost/manifest.yaml"))
	})
	It("should handle real url", func() {
		result, err := adapter.newManifestRequest(
			&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
			"https://github.com/glasskube/glasskube/manifest.yaml",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.URL.String()).To(Equal("https://github.com/glasskube/glasskube/manifest.yaml"))
	})

	Context("with basic auth", func() {
		BeforeEach(func() {
			adapter.repo = fake.ClientsetWithClient(fake.EmptyClientWithAuth(auth.Basic("test", "test")))
		})

		It("should add auth header for relative url", func() {
			result, err := adapter.newManifestRequest(
				&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
				"manifest.yaml",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.URL.String()).To(Equal("http://localhost/packages/foo/manifest.yaml"))
			user, pass, ok := result.BasicAuth()
			Expect(ok).To(BeTrueBecause("basic auth should be added for relative urls"))
			Expect(user).To(Equal("test"))
			Expect(pass).To(Equal("test"))
		})

		It("should NOT add auth header for absolute url", func() {
			result, err := adapter.newManifestRequest(
				&v1alpha1.PackageInfo{Status: v1alpha1.PackageInfoStatus{ResolvedUrl: "http://localhost/packages/foo/package.yaml"}},
				"https://github.com/glasskube/glasskube/manifest.yaml",
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.URL.String()).To(Equal("https://github.com/glasskube/glasskube/manifest.yaml"))
			user, pass, ok := result.BasicAuth()
			Expect(ok).To(BeFalseBecause("basic auth must NOT be added for absolute urls"))
			Expect(user).To(Equal(""))
			Expect(pass).To(Equal(""))
		})
	})
})
