# PagerDuty Short Circuiter
`pdcli` is an integration of go-pagerduty and ocm-container which lets you spawn ocm-container with automatic cluster login and other features based on the pagerDuty alert.

***Note that pdcli is not a reinvention of another pagerDuty CLI tool instead it is a wrapper over go-pagerduty and provides you with all go-pagerduty cli features and much more.***

## Features:

- Users can select the alert, which they want to work upon, through CLI, and they will be able to quickly login to the cluster, using ocm-container, without having to copy-paste information from the alert metadata and will be taken straight to the problematic namespace.
- Users will be provided with alert metadata in the terminal.
- Can query who is oncall.
- `pdcli` requires zero configuration, just one-time login is required.
- Sets helpful environment variables inside ocm-container like $NAMESPACE, $JOB, $POD, $INSTANCE, $VERSION based on the alert metadata.

## Installation

### Option 1: Build from source
First clone the repository somewhere in your $PATH. A common place would be within your $GOPATH.

Example:

```
$ mkdir $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ cd $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ git clone https://github.com/openshift/pagerduty-short-circuiter.git .
```
```
$ make build
```
This command will build the PagerDuty CLI binary, named `pdcli`. This binary will be created in the root of your project folder.

## Running Tests
```
$ make test tests
```
This command runs all the command tests. The test suite uses the [Ginkgo](https://onsi.github.io/ginkgo/) to run comprehensive tests using Behavior-Driven Development.

## Maintainers
- Krishnanunni B (krb@redhat.com)
- Mitali Bhalla (mbhalla@redhat.com)
- Supreeth Basabattini (sbasabat@redhat.com)
- Dominic Finn (dofinn@redhat.com)
- Pooja Rani (prani@redhat.com)

## Feedback
Please help us improve. To contact the SRE-P team:

- @sd-sre-platform in slack channel #sd-sre-platform (CoreOS workspace)
- Or reach out to the maintainers.