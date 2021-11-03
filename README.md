# PagerDuty Short Circuiter
`pdcli` is an integration of [go-pagerduty](https://github.com/PagerDuty/go-pagerduty) and [ocm-container](https://github.com/openshift/ocm-container) which lets you spawn ocm-container with automatic cluster login and other features based on the PagerDuty alert.

***Note that pdcli is not a reinvention of another PagerDuty CLI tool instead it is a wrapper over go-pagerduty and provides you with all go-pagerduty cli features and much more.***

## Features:

- Users can select the alert, which they want to work upon, through CLI, and they will be able to quickly login to the cluster, using ocm-container, without having to copy-paste information from the alert metadata.
- Users will be provided with alert metadata in the terminal.
- Users can switch between different PagerDuty teams they're a part of.
- Users can acknowledge incidents assigned to them.
- Users can query who is oncall for each escalation.
- Users can query when are they scheduled next for oncall.
- `pdcli` requires zero configuration, just one-time login is required.

## Prerequisites:

You will need to have [ocm-container](https://github.com/openshift/ocm-container) installed and configured locally for the cluster login functionality to work.

For further installation instructions please follow the repository [docs](https://github.com/openshift/ocm-container#readme).

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

## Alerts

To view the PagerDuty alerts, use the command:

```
pdcli alerts
```
This will list all the alerts assigned to **self** by default.

You can modify the alerts returned with the `assigned-to` option, you can either choose to list alerts which are assigned to **self**, **team** or **silentTest**.

When viewing alerts assigned to *self*, only acknowledged incident alerts are displayed.

### Alerts View Navigation

By default, all the incident alerts are displayed in the main view.

| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| View resolved alerts                                           | `1`                           | Displays all alerts with status resolved.                              |
| View triggered alerts                                          | `2`                           | Displays all unresolved alerts.                                        |
| View acknowledged incidents                                    | `3`                           | Displays all acknowledged incidents.                                   |
| View triggered incidents                                       | `4`                           | Displays all open incidents assigned to the user.                      |
| View high alerts                                               | `H` / `h`                     | In the trigerred alerts view, when pressed, displays all the alerts with urgency high.|
| View low alerts                                                | `L` / `l`                     | In the trigerred alerts view, when pressed, displays all the alerts with urgency low.|
| View alert data                                                | `Enter`⏎                      | Displays the alert details.                                        |
| Cluster login                                                  | `Y` / `y`                     | In the alert details view, once pressed, spawns an ocm-container instance and proceeds with login into the alert specific cluster.|
| Go back                                                        | `Esc`                         | Navigate to the previous page.                                         |
| Quit                                                           | `Q` / `q`                     | Exit the application.                                                  |


### Incidents View Navigation

When a user navigates to `[3]` trigerred incidents page.

| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| Select/Deselect incident                                       | `Enter`⏎                      | Add or Remove an incident to be acknowledged.                          |
| Acknowledge incident(s)                                        | `ctrl-a`                      | Acknowledge the selected incidents.                                    |
| Go back                                                        | `Esc`                         | Navigate back to alerts main view.                                     |


### View Incident Alerts

To view alerts related to a particular incident, use the command:

```
pdcli alerts <Incident ID>
```
This will list all the alerts belonging to that incident.


### Options
```
--assigned-to          Filter alerts based on user or team (default "self")
--columns              Specify which columns to display separated by commas without any space in between 
                       (default "incident.id,alert,cluster.name,cluster.id,status,severity")
```

## Oncall

To view the current oncalls as per PagerDuty, use the command:

```
pdcli oncall
```
### Oncall View Navigation

By default, all the escalations and Oncalls are displayed for team **Platform-SRE** in the main view.

| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| All Teams Oncall                                               | `A`                           | Displays Escalations and Oncalls for all teams.                              |
| Your Next Oncall Schedule                                      | `N`                           | Displays your Oncall schedule.                                        |
| Go back                                                        | `Esc`                         | Navigate to the previous page.                                         |
| Quit                                                           | `Q` / `q`                     | Exit the application.                                                  |




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
