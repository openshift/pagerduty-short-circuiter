# PagerDuty Short Circuiter
`pdcli` is an integration of [go-pagerduty](https://github.com/PagerDuty/go-pagerduty) and [ocm-container](https://github.com/openshift/ocm-container) which lets you spawn ocm-container with automatic cluster login and other features based on the PagerDuty alert.

***Note that pdcli is not a reinvention of another PagerDuty CLI tool instead it is a wrapper over go-pagerduty and provides you with all go-pagerduty cli features and much more.***

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

Build and copy `pdcli` to your $GOPATH/bin:

```
$ make install
```
### Option 2: Build from source

This command will build the PagerDuty CLI binary, named `pdcli`. This binary will be created in the root of your project folder.

```
$ make build
```
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

## Teams

A user account might belong to a single or multiple pagerduty teams.

To set (or) change your pdcli team, use the command:

```
pdcli teams
```
This will list out all the teams a user is a part of and will prompt the user to select a team for pdcli.

## View Alerts

To view the alerts triggered by PagerDuty, use the command:

```
pdcli alerts
```
This will list all the high alerts assigned to **self** by default.

You can modify the alerts returned with the `assigned-to` option, you can either choose to list alerts which are assigned to *self, team* or *silentTest*.

When you use the option `assigned-to=team`, it will fetch all the alerts assigned to **Platform-SRE** team.

### Interactive Mode

To view alerts in interactive mode, use the command:

```
pdcli alerts -i
```
This will prompt the user to select an incident from a list of incidents.

Once an incident is selected, all the alerts related to that particular incident are listed.

On alert selection, the alert metadata is displayed and the user is asked if they want to proceed with cluster login.

If yes, then an instance of ocm container is spawned with the cluster already logged in.

### View Incident Alerts

To view alerts related to a particular incident, use the command:

```
pdcli alerts <Incident ID>
```
This will list all the alerts belonging to that incident in interactive mode.

### Acknowledge Incidents

To acknowledge an incident assigned to your user account, use the command:

```
pdcli alerts --ack
```
This will list all the incidents currently assigned to self and will prompt the user to choose incident(s) to be acknowledged.

A user can acknowledge a single incident or multiple incidents at once.

To acknowledge all the incidents assgined to your user account, use the command:

```
pdcli alerts --ack-all
```

### Options
```
--ack                  Select and acknowledge incidents assigned to self
--ack-all              Acknowledge all incidents assigned to self
--assigned-to          Filter alerts based on user or team (default "self")
--high                 View all high alerts (default true)
--low                  View all low alerts
--interactive, -i      View alerts in interactive mode and proceed with cluster login
--columns              Specify which columns to display separated by commas without any space in between 
                       (default "incident.id,alert,cluster.name,cluster.id,status,severity")
```

## Oncall

To view the current oncalls as per PagerDuty, use the command:

```
pdcli oncall
```
This will list all the oncalls for each escalation policy of **Platform-SRE** team by default.

You can choose to see the oncalls for all teams on PagerDuty using **All Teams Oncall** on the interactive window.

To view your next oncall schedule as per PagerDuty, use **Your Next Oncall Schedule** on the interactive window.

### Options
```
[A]              View escalations and oncalls for all teams
[N]              View your oncall schedule
```



## Running Tests
The test suite uses the [Ginkgo](https://onsi.github.io/ginkgo/) to run comprehensive tests using Behavior-Driven Development.<br>
The mocking framework used for testing is [gomock](https://github.com/golang/mock).

This command runs all the tests present within the 'tests' folder by default.

```
$ make test
```
Use the mockgen command to generate source code for a mock class given a Go source file containing interfaces to be mocked.

```
$ mockgen -source=foo.go -destination=mock/foo_mock.go
```

## Maintainers
- Dominic Finn (dofinn@redhat.com)
- Mitali Bhalla (mbhalla@redhat.com)
- Pooja Rani (prani@redhat.com)
- Supreeth Basabattini (sbasabat@redhat.com)
- Tomas Dabašinskas (todabasi@redhat.com)


## Feedback
Please help us improve. Contact the Red Hat SRE-P team via:

- @sd-sre-platform in slack channel #sd-sre-platform (CoreOS workspace)
- Or reach out to the [OWNERS](https://github.com/openshift/pagerduty-short-circuiter/blob/main/OWNERS).
