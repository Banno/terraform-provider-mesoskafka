#!/bin/bash

for x in 0 1 2 3 4 5;
do
    curl http://192.168.33.10:7000/api/brokers/stop?id=$x
    curl http://192.168.33.10:7000/api/brokers/remove?id=$x
done

curl http://192.168.33.10:7000/api/brokers/status | jq .
