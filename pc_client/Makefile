build:
	GOBIN=$(shell pwd)/bin go install -mod=vendor .

arm:
	mkdir -p bin
	GOOS=linux GOARCH=arm GOARM=6 go build -mod=vendor -o=bin/powerclient_arm .

lint:
	golangci-lint run -v ./...

proto:
	cd pb&&protoc -I ../../power-monitor/pb -I../../power-monitor/nanopb/generator/proto/ --go_out=. ../../power-monitor/pb/power.proto