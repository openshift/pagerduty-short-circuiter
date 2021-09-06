# PagerDuty Short Circuiter
`pdcli` is an integration of [go-pagerduty](https://github.com/PagerDuty/go-pagerduty) and [ocm-container](https://github.com/openshift/ocm-container) which lets you spawn ocm-container with automatic cluster login and other features based on the pagerDuty alert.

***Note that pdcli is not a reinvention of another pagerDuty CLI tool instead it is a wrapper over go-pagerduty and provides you with all go-pagerduty cli features and much more.***

## Features:

- Users can select the alert, which they want to work upon, through CLI, and they will be able to quickly login to the cluster, using ocm-container, without having to copy-paste information from the alert metadata and will be taken straight to the problematic namespace.
- Users will be provided with alert metadata in the terminal.
- Can query who is oncall.
- `pdcli` requires zero configuration, just one-time login is required.
- Sets helpful environment variables inside ocm-container like $NAMESPACE, $JOB, $POD, $INSTANCE, $VERSION based on the alert metadata.

## Installation
First clone the repository somewhere in your $PATH. A common place would be within your $GOPATH. <br>

Example:

```
$ mkdir $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ cd $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ git clone https://github.com/openshift/pagerduty-short-circuiter.git .
```
You need to have go installed, the minimal version required is go 1.15.

### Option 1: Install binary

```
$ make install
```
This command will fetch, build a binary named `pdcli` and install them to your $GOPATH/bin, you should and move this binary onto your $PATH if desired.

### Option 2: Build from source
```
$ make build
```
This command will build the PagerDuty CLI binary, named `pdcli`. This binary will be created in the root of your project folder.

## Login

- In order to use any PagerDuty CLI features you will have to login first.
- This tool needs zero-configuration, thus only a one-time login is required.
- In order to login, a user must have a valid PagerDuty API key.

To log into PagerDuty CLI use the command:

```
pdcli login
```
This will prompt the user for an API key with necessary instructions of how to generate one. The API key will be saved for future use to the `~/.config/pagerduty-cli/config.json` file.

The `login` command has options to login overwriting the existing API key. For example, if you want to login via another user account or your API key has changed, you can login like this:

```
pdcli login --api-key <api-key>
```
Once logged in you need not login ever again unless there is a change in the API key.

## Running Tests
The test suite uses the [Ginkgo](https://onsi.github.io/ginkgo/) to run comprehensive tests using Behavior-Driven Development.

```
$ make test
```
This command runs all the tests present within the 'tests' folder by default.


## Maintainers
- Krishnanunni B (krb@redhat.com)
- Mitali Bhalla (mbhalla@redhat.com)
- Supreeth Basabattini (sbasabat@redhat.com)
- Dominic Finn (dofinn@redhat.com)
- Pooja Rani (prani@redhat.com)

## Feedback
Please help us improve. Contact the Red Hat SRE-P team via:

- @sd-sre-platform in slack channel #sd-sre-platform (CoreOS workspace)
- Or reach out to the [OWNERS](https://github.com/openshift/pagerduty-short-circuiter/blob/main/OWNERS).