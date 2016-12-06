.PHONY: vet linux osx build test release

vet:
	go tool vet *.go mesoskafka/*.go

linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/terraform-provider-mesoskafka-linux .

osx:
	GOOS=darwin GOARCH=386 go build -o bin/terraform-provider-mesoskafka-osx .

build: vet osx linux
	go install .

test: build
	big inventory add marathon kafka-scheduler
	big up -d marathon kafka-scheduler
	Sleep 5
	TF_ACC=yes MESOS_KAFKA_URL="http://dev.banno.com:7000" go test ./mesoskafka -v

release:
	./bin/release.sh
