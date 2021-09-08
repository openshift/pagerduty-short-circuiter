package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/pagerduty-short-circuiter/pkg/constants"
)

var _ = Describe("pdcli login", func() {

	var jsonTemplate string

	BeforeEach(func() {
		jsonTemplate = `{ 
			"api_key" : "` + constants.SampleKey + `"
			}`

	})

	When("the login command is run", func() {
		It("creates a configuration file", func() {
			result := NewCommand().
				Args(
					"login",
				).
				GetStdIn(constants.SampleKey + "\n").
				Run()

			Expect(result.ConfigFile()).ToNot(BeEmpty())
			Expect(result.ConfigString()).To(MatchJSON(jsonTemplate))
		})
	})

	When("the login command is run with the 'api-key' option", func() {
		It("creates a configuration file", func() {

			result := NewCommand().
				Args(
					"login",
					"--api-key",
					constants.SampleKey,
				).
				Run()

			Expect(result.ConfigFile()).ToNot(BeEmpty())
			Expect(result.ConfigString()).To(MatchJSON(jsonTemplate))
		})
	})

	When("the login command is run", func() {
		It("the sample key is used against PagerDuty REST API", func() {

			result := NewCommand().
				Args(
					"login",
					"--api-key",
					constants.SampleKey,
				).
				Run()

			Expect(result.ConfigFile()).ToNot(BeEmpty())
			Expect(result.ExitCode()).ToNot(BeZero())
			Expect(result.ErrString()).To(ContainSubstring("Unauthorized"))
		})
	})
})
