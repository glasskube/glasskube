package plain

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlainManifestAdapter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PlainManifestAdapter Suite")
}
