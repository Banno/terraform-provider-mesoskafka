provider "mesoskafka" {
  url = "http://dev.banno.com:7000"
}

resource "mesoskafka_cluster" "broker-example" {
   broker_count = 3
   cpus = 0.1
   memory = 256
}
