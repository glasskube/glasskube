package semver

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSemver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Semver Suite")
}
