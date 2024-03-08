<!--
Copyright 2024 Deutsche Telekom IT GmbH

SPDX-License-Identifier: Apache-2.0
-->

<p align="center">
  <h1 align="center">Cosmoparrot</h1>
</p>

<p align="center">
  A simple HTTP based echo server.
</p>

[![REUSE status](https://api.reuse.software/badge/github.com/telekom/pubsub-horizon-cosmoparrot)](https://api.reuse.software/info/github.com/telekom/pubsub-horizon-cosmoparrot)
[![Go Test](https://github.com/telekom/pubsub-horizon-cosmoparrot/actions/workflows/go-test.yml/badge.svg)](https://github.com/telekom/pubsub-horizon-cosmoparrot/actions/workflows/go-test.yml)

## Overview
Cosmosparrot simple HTTP based echo server designed to provide a response that mirrors the contents included in the initial request.  
It was initially created for Pub/Sub end-to-end test scenarios where it is important to simulate an event message consumer that responds to HTTP (callback) requests.

## Running Cosmoparrot
To run Cosmoparrot run the `cosmoparrot` executable to start the server.

Alternatively you can use the provided Dockerfile to build and run the server in a container.  

Build the image:
```bash
docker build -t cosmoparrot .
```  
Run the container:
```bash
docker run -p 8080:8080 cosmoparrot
```

## Configuration
Cosmoparrot supports configuration via environment variables and/or a configuration file (`config.yml`). The configuration file has to be located in the same directory as the executable.


| Path                       | Variable                 | Type | Default | Description                                                                              |
|----------------------------|--------------------------|------|---------|------------------------------------------------------------------------------------------|
| port                       | COSMOPARROT_PORT         | int  | 8080    | Sets the port to listen on.                                                              |
| responseCode               | COSMOPARROT_RESPONSECODE | int  | 200     | Enforces a specific HTTP response code. Can be used to test different consumer behavior. |

## Code of Conduct
This project has adopted the [Contributor Covenant](https://www.contributor-covenant.org/) in version 2.1 as our code of conduct. Please see the details in our [Code of Conduct](CODE_OF_CONDUCT.md). All contributors must abide by the code of conduct.
By participating in this project, you agree to abide by its [Code of Conduct](./CODE_OF_CONDUCT.md) at all times.

## Licensing
This project follows the [REUSE standard for software licensing](https://reuse.software/).
Each file contains copyright and license information, and license texts can be found in the [LICENSES](./LICENSES) directory.

For more information visit https://reuse.software/.
For a comprehensive guide on how to use REUSE for licensing in this repository, visit https://telekom.github.io/reuse-template/.   
A brief summary follows below:

The [reuse tool](https://github.com/fsfe/reuse-tool) can be used to verify and establish compliance when new files are added.

For more information on the reuse tool visit https://github.com/fsfe/reuse-tool.

