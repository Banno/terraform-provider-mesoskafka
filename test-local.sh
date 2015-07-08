#!/bin/bash

TF_ACC=yes MESOS_KAFKA_URL="http://dev.banno.com:7000" go test ./mesoskafka -v
