#!/bin/bash

curl http://192.168.33.10:7000/api/brokers/stop?id=0
curl http://192.168.33.10:7000/api/brokers/stop?id=1
curl http://192.168.33.10:7000/api/brokers/stop?id=2

curl http://192.168.33.10:7000/api/brokers/remove?id=0
curl http://192.168.33.10:7000/api/brokers/remove?id=1
curl http://192.168.33.10:7000/api/brokers/remove?id=2


curl http://192.168.33.10:7000/api/brokers/status | jq .
