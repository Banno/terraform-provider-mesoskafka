#Mesos Kafka Provider For Terraform

Allows you to define Kafka brokers that run inside Mesos using the
(Mesos Kafka Framework)[https://github.com/mesos/kafka]

##Requirements: 
* A working Mesos cluster (with the mesos executor enabled)
* Mesos slaves running Java 8
* A running (Mesos Kafka Framework Scheduler)[https://github.com/mesos/kafka]

## Provider

The provider can be configured manually via
```
provider "mesoskafka" {
  url = "http://kafka-scheduler:7000"
}
```

or with environmental variables:
```
export MESOS_KAFKA_URL=http://kafka-scheduler:7000
```

##Minimal configuration
```
resource "mesoskafka_cluster" "my_cluster" {
  broker_count = 3
  cpus = 1
  memory = 256
}
```

##How to build
```bash
$ go get
$ go install
```
