# Dozer

## Morpheus Processes with Webhooks

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/spoonboy-io/dozer?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/spoonboy-io/dozer?style=flat-square)](https://goreportcard.com/report/github.com/spoonboy-io/dozer)
[![DeepSource](https://deepsource.io/gh/spoonboy-io/dozer.svg/?label=active+issues&token=uYY_4Kwjq9MnjT7TzykEyv-J)](https://deepsource.io/gh/spoonboy-io/dozer/?ref=repository-badge)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/spoonboy-io/dozer/Build?style=flat-square)](https://github.com/spoonboy-io/dozer/actions/workflows/build.yml)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/spoonboy-io/dozer/Unit%20Test/master?label=tests&style=flat-square)](https://github.com/spoonboy-io/dozer/actions/workflows/unit_test.yml)

[![GitHub Release Date](https://img.shields.io/github/release-date/spoonboy-io/dozer?style=flat-square)](https://github.com/spoonboy-io/dozer/releases)
[![GitHub commits since latest release (by date)](https://img.shields.io/github/commits-since/spoonboy-io/dozer/latest?style=flat-square)](https://github.com/spoonboy-io/dozer/commits)
[![GitHub](https://img.shields.io/github/license/spoonboy-io/dozer?label=license&style=flat-square)](LICENSE)

## About

Dozer watches [Morpheus CMP](https://morpheusdata.com) processes/events. 
It will notify external applications of Morpheus events
via HTTP request (webhook) based on YAML configuration you specify.

## Releases

You can find the [latest software here](https://github.com/spoonboy-io/dozer/releases/latest).

### Get Started

Dozer polls the Morpheus database so needs credentials. The `morpheus` user can be used, but it is preferable to 
create an additional user with SELECT privileges on the `process` and `process_type` tables.

A `mysql.env` file should be created in the same directory as the application from which the database user configuration
will be read. The following example shows the environment variables used by Dozer which should be included in `mysql.env`:

```bash
## MySQL Config
MYSQL_SERVER=127.0.0.1
MYSQL_USER=dozer
MYSQL_PASSWORD=xxxxa8aca0de5dab5fa1bxxxxx

## Optional to override default database poll interval (5 seconds)
POLL_INTERVAL_SECONDS=3
```

### Webhook Configuration

Webhooks and their triggers are configured in a YAML file `webhook.yaml` which should reside in the same directory
as the Dozer application. An example configuration, showing a single Webhook is shown below:

```YAML
---
- webhook:
    description: Hook example with trigger that runs when status is `complete`
    url: https://webhook-endpoint.com
    method: POST
    requestBody: |
        {
            "id": {{.Id}},
            "updatedBy": "{{.UpdatedBy}}",
            "status": "{{.Status}}"
        }
    token: BEARER xxxxxxxxxxxxxx
    triggers:
      status: complete
```
GET and POST methods are supported. If POST method Dozer will look for `requestBody`.

If `token` is supplied it will be sent in the AUTHORIZATION header.

Variables which contain information about the Morpheus process can be interpolated in the `requestBody` using the standard Golang
templating format. A complete list can be found [here](https://github.com/spoonboy-io/dozer/blob/master/internal/hook/send.go#L15).

### Triggers

Currently, Webhook triggers can be specified on `status`, `processType`, `taskName`, `accountId` and `createdBy`. They are 
evaluated on processes which have finished running, not in-progress processes. Triggers are additive - all conditions must 
be satisfied for the Webhook to fire.

| Trigger 	        | Description 	                                            | YAML Example                  |
|---------	        |-------------	                                            | ---------	                    |
| `status`          | Runs when the process is complete or failed       	    | `status: failed`              |
| `processType`     | Runs for a specific process type ([see list here](https://github.com/spoonboy-io/dozer/blob/master/internal/morpheus/processType.go#L11))       | `processType: localWorkflow`  |
| `taskName`        | Runs for a given task name            	                | `taskName: Hello World`       |
| `accountId`       | Runs for specific tenant id           	                | `accountId: 2`        	    |
| `createdBy`       | Runs for processes created by a specific user            	| `createdBy: admin`        	|


### Installation
Grab the tar.gz or zip archive for your OS from the [releases page](https://github.com/spoonboy-io/dozer/releases/latest).

Unpack it to the target host, and then start the server.

```
./dozer
```

Or with nohup..

```
nohup ./dozer &
```

### Development Opportunities

- Add more triggers such as `zoneId`, `instanceName`, `appName`, `containerName`
- Retry and blacklisting for webhooks that fail
- Other notification mechanisms such as email or messaging protocol
- Run as a service

### License
Licensed under [Mozilla Public License 2.0](LICENSE)
