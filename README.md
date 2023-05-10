# Kite
`kite` is an integration of [go-pagerduty](https://github.com/PagerDuty/go-pagerduty) and [ocm-container](https://github.com/openshift/ocm-container) which lets you spawn ocm-container with automatic cluster login and other features based on the PagerDuty alert.

***Note that kite is not a reinvention of another PagerDuty CLI tool instead it is a wrapper over go-pagerduty and provides you with all go-pagerduty cli features and much more.***

## Features

- Users can select the alert, which they want to work upon, through CLI, and they will be able to quickly login to the cluster, using ocm-container, without having to copy-paste information from the alert metadata.
- Users will be provided with alert metadata in the terminal.
- Users can switch between different PagerDuty teams they're a part of.
- Users can acknowledge incidents assigned to them.
- Users can see whom the alerts are being assigned to
- Users can navigate between windows and tabs using key shortcuts
- Users can query who is oncall for each escalation.
- Users can query when are they scheduled next for oncall.
- Users can navigate between previous and next layer of oncall schedule.
- Users can view the SOP of a particular alert 
- `kite` requires zero configuration, just one-time login is required.

## Prerequisites

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

Build and copy `kite` to your $GOPATH/bin:

```
$ make install
```
### Option 2: Build from source

This command will build the PagerDuty CLI binary, named `kite`. This binary will be created in the root of your project folder.

```
$ make build
```
## Login

- In order to use any PagerDuty CLI features you will have to login first.
- This tool needs zero-configuration, thus only a one-time login is required.
- In order to login, a user must have a valid PagerDuty API key.

To log into PagerDuty CLI use the command:

```
kite login
```
This will prompt the user for an API key with necessary instructions of how to generate one. The API key will be saved for future use to the `~/.config/kite/config.json` file.

The `login` command has options to login overwriting the existing API key. For example, if you want to login via another user account or your API key has changed, you can login like this:

```
kite login --api-key <api-key>
```
Upon login, kite will ask the user to enter a [Github Access Token](https://github.com/settings/tokens). Users need to create a classic access token with read-only access.

## Teams

A user account might belong to a single or multiple pagerduty teams.

To set (or) change your kite team, use the command:

```
kite teams
```
This will list out all the teams a user is a part of and will prompt the user to select a team for kite.

## Alerts

To view the PagerDuty alerts, use the command:

```
kite alerts
```
To view the PagerDuty alerts that have been assigned to your team, use the command:

```
kite alerts --assigned-to team
```
To view the PagerDuty alerts that have been assigned to silentTest, use the command:

```
kite alerts --assigned-to silentTest
```
This will list all the alerts assigned to **self** by default.

You can modify the alerts returned with the `assigned-to` option, you can either choose to list alerts which are assigned to **self**, **team** or **silentTest**.

When viewing alerts assigned to *self*, only acknowledged incident alerts are displayed.

### Alerts View Navigation

By default, all the incident alerts are displayed in the main view.

| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| View triggered alerts                                          | `1`                           | Displays all unresolved alerts.                                        |
| View acknowledged incidents                                    | `2`                           | Displays all acknowledged incidents.                                   |
| View triggered incidents                                       | `3`                           | Displays all open incidents assigned to the user.                      |
| View alert data                                                | `Enter`⏎                      | Displays the alert details.                                            |
| Cluster login                                                  | `Y` / `y`                     | In the alert details view, once pressed, spawns an ocm-container instance and proceeds with login into the alert specific cluster.|
| View SOP                                                       | `S` / `s`                     | Displays the SOP for that alert                                        |
| View Service Logs                                              | `L` / `l`                     | Displays the Service Logs                                              |
| Refresh Alerts                                                 | `R` / `r`                     | Refreshes the alerts                                                   |
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
kite alerts <Incident ID>
```
This will list all the alerts belonging to that incident.

### View Service Logs
An alerting cluster's service logs can be viewed while viewing the alert data by pressing `L/l`.

### Options
```
--assigned-to          Filter alerts based on user or team (default "self") 
--columns              Specify which columns to display separated by commas without any space in between 
                       (default "incident.id,alert,cluster.name,cluster.id,status,severity")
```

## Oncall

To view the current oncalls as per PagerDuty, use the command:

```
kite oncall
```
### Oncall View Navigation

By default, all the escalations and Oncalls are displayed for team **Platform-SRE** in the main view.

| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| All Teams Oncall                                               | `A` / `a`                     | Displays Escalations and Oncalls for all teams.                        |
| Your Next Oncall Schedule                                      | `N` / `n`                     | Displays your Oncall schedule.                                         |
| Previous Layer Oncall                                          | `[<-]`                        | Displays previous layer of Oncall schedule.                            |
| Next Layer Oncall                                              | `[->]`                        | Displays next layer Oncall schedule.                                   |
| Go back                                                        | `Esc`                         | Navigate to the previous page.                                         |
| Quit                                                           | `Q` / `q`                     | Exit the application.                                                  |

## Terminal Multiplexer

The terminal multiplexer is integrated with both alerts and oncall view. New Shell Prompt is opened up when we add a tab. You can also navigate between windows and tabs using key shortcuts.

### Terminal Multiplexer Navigation


| Action                                                         | Key                           | Comment                                                                |
|----------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------|
| Next Slide                                                     | Ctrl + `N` / `n`              | Moves to next slide                                                    |
| Previous Slide                                                 | Ctrl + `P` / `p`              | Moves to previous slide                                                |
| Add Slide                                                      | Ctrl + `A` / `a`              | Adds a new slide                                                       |
| Add Slide `ocm-container`                                      | Ctrl + `O` / `o`              | Adds a new slide (ocm-container)                                       |
| Change to Slide with [Num]                                     | Ctrl + `B` + [Num]            | Move to a particular slide                                             |
| Exit Slide                                                     | Ctrl + `E` / `e`              | Exit a slide                                                           |
| Quit                                                           | Ctrl + `Q` / `q`              | Exit the application.                                                  |

## In-built SOP

To view SOP related to a particular alert a user can run these commands:
* To show alerts use the command :
```
kite alerts
```
* To select an alert and view its metadata press `enter`
* To view SOP, press `S`

## Terminal
To Select a preferred terminal emulator, use the command:

```
kite terminal
```
This will list all the terminal emulators supported by the system and will prompt the user to select an emulator for kite.
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
## List of known Bugs
There are few bugs that have crawled up during the development. These are listed below :
- exit command crashing kite in certain instances when used in shell
- Auto completion changes moves cursor to different position but text typing continues from same position
- Bash initialization is slow in certain cases
- Use of mouse causes random text to be typed in the terminal
- Block cursor does not move with the text being typed (Quick Fix : Changing cursor style to line)
- When viewing an SOP, if the relative link is present for another SOP file, the new SOP might not get opened.(The URL needs to be absolute for the SOP navigtaion to work

## Maintainers
- Mitali Bhalla (mbhalla@redhat.com)
- Supreeth Basabattini (sbasabat@redhat.com)
- Tomas Dabašinskas (todabasi@redhat.com)

## List 

## Feedback
Please help us improve. Contact the Red Hat SRE-P team via:

- @sd-sre-platform in slack channel #sd-sre-platform (CoreOS workspace)
- Or reach out to the [OWNERS](https://github.com/openshift/pagerduty-short-circuiter/blob/main/OWNERS).
