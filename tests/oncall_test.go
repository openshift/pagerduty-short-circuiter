package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
)

var _ = Describe("pdcli oncall", func() {

	When("Configuration file doesn't exist", func() {
		It("throws an error", func() {
			cmd := NewCommand().Args("oncall").Run()
			Expect(cmd.ExitCode()).NotTo(BeZero())
			Expect(cmd.ErrString()).ToNot(BeEmpty())
		})

	})
})
