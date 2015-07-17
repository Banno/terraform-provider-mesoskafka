#Mesos Kafka Provider For Terraform

Allows you to define Kafka brokers that run inside Mesos using the
[Mesos Kafka Framework](https://github.com/mesos/kafka)

##Requirements: 
* A working Mesos cluster (with the mesos executor enabled)
* Mesos slaves running Java 8
* A running [Mesos Kafka Framework](https://github.com/mesos/kafka)

## Provider

The provider can be configured manually via
```
provider "mesoskafka" {
  url = "http://mesoskafka:7000"
}
```

or with environmental variables:
```
export MESOS_KAFKA_URL=http://mesoskafka:7000
```

##Minimal configuration
```
resource "mesoskafka_cluster" "my_cluster" {
  broker_count = 3
  cpus = 1
  memory = 256
}
```

##Full configuration
```
resource "mesoskafka_cluster" "my_cluster" {
  broker_count = 3
  cpus = 1
  memory = 256
  heap = 128
  jvm_options = "-Xms128m"
  options = "file:server.properties"
  logfourj_options = "file:log4j.properties"
  failover_delay = "14s"
  failover_max_delay = "5s"
  failover_max_tries = 5
}
```

##How to install
```bash
$ go get github.com/Banno/terraform-provider-mesoskafka
```
