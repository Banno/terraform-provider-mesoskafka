#!/bin/bash

TF_ACC=yes MESOS_KAFKA_URL="http://mesos-master0-aws.ngrayson.dev.banno-internal.com:7000" go test ./mesoskafka -v
