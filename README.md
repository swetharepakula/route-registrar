route-registrar
===============

A standalone executable written in golang that continuously broadcasts a routes to the [gorouter](https://github.com/cloudfoundry/gorouter).  This is designed to be a general purpose solution, packaged as a BOSH job to be colocated with components that need to broadcast their routes to the gorouter, so that those components don't need to maintain logic for route registration.

* CI: [Concourse](https://cf-routing.ci.cf-app.com/pipelines/route-registrar)

## Usage

1. Run the following command to install route-registrar
  ```
  go get github.com/cloudfoundry-incubator/route-registrar
  ```

1. The route-registrar expects a configuration YAML file like the one below:
  ```yaml
  message_bus_servers:
  - host: NATS_SERVER_HOST
    user: NATS_SERVER_USERNAME
    password: NATS_SERVER_PASSWORD
  host: HOSTNAME_OR_IP_OF_ROUTE_DESTINATION
  routes:
  - name: SOME_ROUTE_NAME
    port: PORT_OF_ROUTE_DESTINATION
    tags:
      optional_tag_field: some_tag_value
      another_tag_field: some_other_value
    uris:
    - some_source_uri_for_the_router_to_map_to_the_destination
    - some_other_source_uri_for_the_router_to_map_to_the_destination
    route_service_url: https://route-service.example.com
    registration_interval: REGISTRATION_INTERVAL # required
    health_check: # optional
      name: HEALTH_CHECK_NAME
      script_path: /path/to/check/executable
      timeout: HEALTH_CHECK_TIMEOUT # optional
  ```
  - `message_bus_servers` is an array of data with location and credentials for the NATS servers; route-registrar currently registers and deregisters routes via NATS messages.
  - `host` is the destination hostname or IP for the routes being registered.
  - for each route collection, `name` must be provided and be a string.
  - for each route collection, `port` must be provided and must be a positive integer > 1.
  - for each route collection, `uris` must be provided and be a non empty array of strings.  All URIs in a given route collection will be mapped to the same host and port.
  - for each route collection, `registration_interval` must be provided and be a string with units (e.g. "20s"). It must parse to a positive time duration e.g. "-5s" is not permitted.
  - for each route collection, `route_service_url` is optional and enables the component to register a route service for that route.  
  - for each route collection, `health_check` is optional and explained in more detail below.

1. Run route-registrar binaries using the following command
  ```
  ./bin/route-registrar -configPath=FILE_PATH_TO_CONFIG_YML --pidFile=PATH_TO_PIDFILE
  ```

### Health check

If the `health_check` is not configured for a route collection, the routes are continually registered according to the `registration_interval`.

If the `health_check` is configured, the executable provided at `health_check.script_path` is invoked and the following applies:
- if the executable exits with success, the routes are registered.
- if the executable exits with error, the routes are deregistered.
- if `health_check.timeout` is configured, it must parse to a positive time duration (similar to `registration_interval`), and the executable must exit within the timeout. If the executable does not terminate within the timeout, it is forcibly terminated (with `SIGKILL`) and the routes are deregistered.
- if `health_check.timeout` is not configured, the executable must exit within half the `registration_interval`. If the executable does not terminate within the timeout, it is forcibly terminated (with `SIGKILL`) and the routes are deregistered.

## BOSH release

This program is packaged as a [job](https://github.com/cloudfoundry/cf-release/tree/master/jobs/route_registrar) and a [package](https://github.com/cloudfoundry/cf-release/tree/master/packages/route_registrar) in the [cf-release](https://github.com/cloudfoundry/cf-release)
BOSH release, it can be colocated with the following manifest changes:

```yaml
releases:
- name: cf
- name: my-release

jobs:
- name: myJob
  templates:
  - name: my-job
    release: my-release
  - name: route_registrar
    release: cf
  properties:
    route_registrar:
      # ...
      # [ see bosh job spec ]

```

## Development

### Dependencies

Dependencies are saved using [Godep](https://github.com/tools/godep) with `godep save -r` (import path re-writing).
Just clone the repo to your `GOPATH` and the dependencies should be available.

### Running tests

1. Install the ginkgo binary with `go get`:
  ```
  go get github.com/onsi/ginkgo/ginkgo
  ```

1. Run tests, by running the following command from root of this repository
  ```
  bin/test
  ```
