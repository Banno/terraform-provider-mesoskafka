#!/bin/bash

TF_ACC=yes MESOS_KAFKA_URL="http://192.168.33.10:7000" go test ./mesoskafka -v
