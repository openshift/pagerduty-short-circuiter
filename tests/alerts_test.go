package tests

import (
	
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("view alerts test", func() {

	When("configuration file doesn't exist", func() {
		It("throws an error", func() {
			cmd := NewCommand().Args("alerts").Run()

			Expect(cmd.ExitCode()).NotTo(BeZero())
			Expect(cmd.ErrString()).ToNot(BeEmpty())
		})
	})

})
