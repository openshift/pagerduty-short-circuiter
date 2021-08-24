# PagerDuty Short Circuiter
This project contains PagerDuty CLI to reduce the time taken from when SRE receives a PD alert to when troubleshooting on the cluster actually begins.

- The document of the project design can be found [here](https://docs.google.com/document/d/1VV3bN3WBI-DrJ59jOciA5tyvnxhldSVW-lDGNu2ONiw/edit?usp=sharing).
- Jira Epics for this project can be found [here](https://issues.redhat.com/browse/OSD-8102).

## Installation

### Option 1: Build from source
First clone the repository somewhere in your $PATH. A common place would be within your $GOPATH.

Example:

```
$ mkdir $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ cd $GOPATH/src/github.com/openshift/pagerduty-short-circuiter
$ git clone https://github.com/openshift/pagerduty-short-circuiter.git
```
```
$ make build
```
This command will build the PagerDuty CLI binary, named `pdcli`. This binary will be placed in cmd/pdcli folder of your project directory.

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