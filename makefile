.PHONY: vet linux osx build test release

vet:
	go tool vet *.go mesoskafka/*.go

linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/terraform-provider-mesoskafka-linux .

osx:
	GOOS=darwin GOARCH=386 go build -o bin/terraform-provider-mesoskafka-osx .

install:
	go install .

os=$(shell uname)
ifeq ($(os),Darwin)
  DEV_BANNO_IP=192.168.99.100
  LOCAL_BANNO_IP=192.168.99.1
else
  DEV_BANNO_IP=127.0.0.1
  LOCAL_BANNO_IP=127.0.0.1
endif

test: install
	DEV_BANNO_IP=$(DEV_BANNO_IP) LOCAL_BANNO_IP=$(LOCAL_BANNO_IP) docker-compose up -d
	sleep 5
	TF_ACC=yes MESOS_KAFKA_URL="http://dev.banno.com:7000" go test -timeout 20m ./mesoskafka -v
	docker-compose kill
	docker-compose rm -f

release: vet linux osx
	./bin/release.sh
