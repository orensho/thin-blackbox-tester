# thin Blackbox Tester #

[![Release](https://img.shields.io/github/release/orensho/thin-blackbox-tester/all.svg)](https://github.com/orensho/thin-blackbox-tester/latest)
[![PkgGoDev](https://pkg.go.dev/badge/orensho/thin-blackbox-testert/)](https://github.com/orensho/thin-blackbox-tester/)

I created this thin Blackbox tester to demonstrate how easy it is to build and automate your ASM testing

## Description

A thin Blackbox tester to fork from for your custom blackbox tester

## Required environment variables

|Env-var                            | Required | Default Value | Description |
|-----------------------------------|--------- |---------------|-------------|
|**TESTER_CONFIG_FILENAME**|yes|/config.yml|location of the exported config file|
|**TESTER_CONFIG_FOLDER**|yes|configuration|folder of the exported config file|
|**TESTER_SHOW_DEBUG_BROWSER**|no|false||
|**TESTER_ENVIRONMENT**|yes|dev||
|**SERVER_LOCAL_LISTEN_IP**|yes|127.0.0.1||
|**SERVER_LOCAL_LISTEN_PORT**|yes|8080||
|**SERVER_SHUTDOWN_GRACE_PERIOD**|yes|10s||
|**METRICS_ENABLE**|yes|true||
|**METRICS_ENVIRONMENT**|yes|dev||
|**METRICS_PREFIX**|yes|thin_blackbox||
|**METRICS_PORT**|yes|8888||


## Endpoints

* METRICS_PORT/metrics

## Build and Run

* Local: just run package 'github.com/orensho/thin-blackbox-tester/service/cmd/service' and default configuration will be used
* Makefile: call ``makefile build_run_darwin``

## Metrics
Calling ``curl SERVER_LOCAL_LISTEN_IP:METRICS_PORT/metrics`` will return the blackbox tester current metrics

## Deployment

Your containerized blackbox tester should be deployed on a workload to provide availability<br />
It is recommended to create a CI pipeline to automate deployment 
