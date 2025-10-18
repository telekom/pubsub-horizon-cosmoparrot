<!--
Copyright 2024 Deutsche Telekom IT GmbH

SPDX-License-Identifier: Apache-2.0
-->

<p align="center">
  <img src="docs/img/cosmoparrot-logo.svg" alt="Cosmoparrot logo" width="200">
  <h1 align="center">Cosmoparrot</h1>
</p>

<p align="center">
  A simple HTTP based echo server.
</p>

<p align="center">
  <a href="#building-cosmoparrot">Building Cosmoparrot</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#running-cosmoparrot">Running Cosmoparrot</a>
</p>

<!--
[![REUSE status](https://api.reuse.software/badge/github.com/telekom/pubsub-horizon-cosmoparrot)](https://api.reuse.software/info/github.com/telekom/pubsub-horizon-cosmoparrot)
-->
[![Go Test](https://github.com/telekom/pubsub-horizon-cosmoparrot/actions/workflows/go-test.yml/badge.svg)](https://github.com/telekom/pubsub-horizon-cosmoparrot/actions/workflows/go-test.yml)

## Overview
Cosmosparrot simple HTTP based echo server designed to provide a response that mirrors the contents included in the initial request.  
It was initially created for Pub/Sub end-to-end test scenarios where it is important to simulate an event message consumer that responds to HTTP (callback) requests.

## Building Cosmoparrot

### Go build

Assuming you have already installed [go](https://go.dev/), simply run the follwoing to build the executable:
```bash
go build
```

> Alternatively, you can also follow the Docker build in the following section if you want to build a Docker image without the need to have Golang installed locally.

### Docker build

This repository provides a multi-stage Dockerfile that will also take care about compiling the software, as well as dockerizing Cosmoparrot. Simply run:

```bash
docker build -t cosmoparrot:latest  . 
```

## Configuration
Cosmoparrot supports configuration via environment variables and/or a configuration file (`config.yml`). The configuration file has to be located in the same directory as the executable.

| Path                        | Variable                              | Type   | Default | Description                                                                              |
|-----------------------------|---------------------------------------|--------|---------|------------------------------------------------------------------------------------------|
| port                        | COSMOPARROT_PORT                      | int    | 8080    | Sets the port to listen on.                                                              |
| responseCode                | COSMOPARROT_RESPONSECODE              | int    | 200     | Enforces a specific HTTP response code. Can be used to test different consumer behavior. |
| methodResponseCodeMapping   | COSMOPARROT_METHODRESPONSECODEMAPPING | string | ""      | Control the HTTP response code per HTTP method, for example: "POST:401"                  |

Additional behavior:

- Query-parameter override: You can override the response code per-request using the `responseCode` query parameter. The server looks for query parameter names in a case-insensitive and tolerant way (for example `responseCode`, `response_code`, `RESPONSECODE`, `Response-Code` are all accepted). When present and valid (an integer between 100 and 599), the query parameter takes precedence over any configured method mapping or global default.

- Method mapping fallback: If no valid `responseCode` query parameter is present, Cosmoparrot checks `methodResponseCodeMapping` for an entry matching the HTTP method (e.g. `GET:202`) and uses that value. If none match, the `responseCode` configuration value is used.

## Running Cosmoparrot
### Locally

Simply run the built `cosmoparrot` executable to start the server:
```shell
./cosmoparrot
```

Alternatively you can run the server in a container: 

```bash
docker run -p 8080:8080 cosmoparrot
```

Notes about background workers:

- Background workers (used for store-write batching and other background tasks) do not start automatically when you only construct the app via `NewApp()` in tests or other programmatic contexts. This keeps tests deterministic and avoids extra goroutines during unit runs.

- When you run the real server (via `Listen()`), Cosmoparrot will explicitly call `StartBackgroundWorkers()` before listening. If you embed or run the app programmatically and need background workers, call `api.StartBackgroundWorkers()` once during startup.

## Deployment

For the deployment of Cosmoparrot you can use Kubernetes deployment `manifest/deployment.yaml` and adjust it to your
needs, or you can use and customize the Heln chart located in `manifest/helm`.

*Helm example:*
```
helm install cosmoparrot ./manifest/helm/cosmoparrot \
  --namespace custom-namespace --create-namespace \
  --set cosmoparrot.storeKeyRequestHeaders="{X-Request-ID,X-Correlation-ID}" \
  --set image.repository=myregistry.com/cosmoparrot \
  --set image.tag=latest \
  --set ingress.enabled=true \
  --set ingress.host=cosmoparrot.mycompany.com \
  --set imagePullSecrets[0].name=my-pull-secret
```

## Contributing

We're committed to open source, so we welcome and encourage everyone to join its developer community and contribute, whether it's through code or feedback.  
By participating in this project, you agree to abide by its [Code of Conduct](./CODE_OF_CONDUCT.md) at all times.

## Code of Conduct
This project has adopted the [Contributor Covenant](https://www.contributor-covenant.org/) in version 2.1 as our code of conduct. Please see the details in our [Code of Conduct](CODE_OF_CONDUCT.md). All contributors must abide by the code of conduct.
By participating in this project, you agree to abide by its [Code of Conduct](./CODE_OF_CONDUCT.md) at all times.

## Licensing

This project follows the [REUSE standard for software licensing](https://reuse.software/). You can find a guide for developers at https://telekom.github.io/reuse-template/.   
Each file contains copyright and license information, and license texts can be found in the [./LICENSES](./LICENSES) folder. For more information visit https://reuse.software/.
