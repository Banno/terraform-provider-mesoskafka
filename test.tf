provider "mesoskafka" {
  url = "http://mesos-master.vagrant:7000"
}

resource "mesoskafka_cluster" "broker-example" {
   broker_count = 5
}
