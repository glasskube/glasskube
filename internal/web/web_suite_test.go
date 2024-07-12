package web

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWeb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Suite")
}
