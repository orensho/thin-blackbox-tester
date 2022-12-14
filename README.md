# Fireglass Blackbox Tester #

Fireglass Blackbox Tester is a Golang server that continuously test isolate and bypass navigation for all production tenants and reports results as metrics 


## Run 

* Local: just run package 'github.com/orensho/thin-slack-blackbox-tester/service/cmd/service' and default configuration will be used
* Kubernetes: wrap in deployment and supply the relevant environment variables

## Environment variables

|Env-var                            | Required | Default Value | Description |
|-----------------------------------|--------- |---------------|-------------|
|**TESTER_CONFIG_FILENAME**|yes|/config.yml|location of the exported config file|
|**TESTER_CONFIG_FOLDER**|yes|configuration|folder of the exported config file|
|**TESTER_SHOW_DEBUG_BROWSER**|no|false||
|**TESTER_ENVIRONMENT**|yes|dev||
|**TESTER_INSTANCE_NAME**|yes|dev|the tester shard instance|
|**SERVER_LOCAL_LISTEN_IP**|yes|127.0.0.1||
|**SERVER_LOCAL_LISTEN_PORT**|yes|8080||
|**SERVER_SHUTDOWN_GRACE_PERIOD**|yes|10s||
|**METRICS_ENABLE**|yes|true||
|**METRICS_ENVIRONMENT**|yes|local||
|**METRICS_PREFIX**|yes|fg_blackbox||
|**METRICS_PORT**|yes|8888||


## Endpoints

* /info
* /health
* /shutdown
* METRICS_PORT/metrics

## Scale out

To support hundreds and more of gateways to test, we will deploy BBOX as a statefulset with 2 replicas
Each BBOX instance will perform hashedGatewayName % TESTER_INSTANCE_NAME ==0 to decide which gateways it should test


## Backlog

- [x] Refactor to BBOX standards
- [x] Review refactor PR
- [x] Review flow
- [x] Review navigate-step
- [x] Integrate with new system rules
- [x] Implement proxies-step
- [x] Implement unit testing
- [x] Implement TF for GCS bucket backed dummy website for isolation and pipeline
- [x] Design build pipeline
- [x] Implement build pipeline
- [x] Design deploy pipeline
- [x] Implement deploy pipeline
- [x] Implement monitoring
