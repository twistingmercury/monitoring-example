# Observability middleware for gin-gonic/gin: Examples

This repository contains an example of using the package [github.com/twistingmercury/monitoring](https://github.com/twistingmercury/monitoring)


## Prerequisites

* Go 1.21.x
* Docker and Docker Compose
* Python v3.x
* The latest [Datadog Agent](https://docs.datadoghq.com/agent/versions/upgrade_to_agent_v7/?tab=linux) Docker image
* A [Datadog](https://www.datadoghq.com/) account

## About

The intent of this examples is to roughly simulate running the DD agent as a side car in a Kubernetes pod, without the need of using minikube locally.

- The [OTel Collector](https://github.com/open-telemetry/opentelemetry-collector) could be used in place of DD agent for traces.
- Prometheus/Grafana could be used to scrape metrics instead of the DD agent.
- Any other tool, like Vector, FluentBit, etc., that can read stdout from a container could be used for logs instead of the DD agent.

This example creates a [docker image](./dockerfile) using Alpine, via a multistage build. 

* The Datadog Agent is configured to expose OTel ports in the [docker-compose.yaml](./docker-compose.yaml#24) file.
* The [volume mapping](/docker-compose.yaml#35) enables the agent to read the logs from the monex container's stout.
* In dockerfile, the following labels permit the agent to scrape prometheus metrics: 
  * [com.datadoghq.ad.check_names='["openmetrics"]'](./dockerfile#18)
  * [com.datadoghq.ad.init_configs='[{}]'](./dockerfile#19)
  * [com.datadoghq.ad.instances='[{"openmetrics_endpoint":"http://monex:9090/metrics","namespace":"example","metrics":["example*"]}]'](./dockerfile#20)

The intent of this examples is to roughly simulate running the DD Agent as a side car.

## Running the example

`make run`, or alternatively `docker compose up`


A python script is provided so that the endpoint is called continuouly to exercise the example service here:[testclient.py](testclient/client.py).

`make client`, or alternatively `python3 ./testclient/client.py`