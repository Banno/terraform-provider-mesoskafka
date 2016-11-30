.PHONY: vet linux osx build test release

vet:
	go tool vet *.go mesoskafka/*.go

linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/terraform-provider-mesoskafka-linux .

osx:
	GOOS=darwin GOARCH=386 go build -o bin/terraform-provider-mesoskafka-osx .

install:
	go install .

test: install
	big destroy
	big stack reset
	big inventory add kafka-scheduler service-registry-watcher
	big up -d
	sleep 5
	TF_ACC=yes MESOS_KAFKA_URL="http://dev.banno.com:7000" go test -timeout 20m ./mesoskafka -v

release: vet linux osx
	./bin/release.sh
