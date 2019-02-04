SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

build: alerts resources gateway

proto:
	protoc -I$(GOPATH)/src -Ipkg/protobuf pkg/protobuf/boolean.proto --go_out=plugins=grpc:pkg/protobuf
	protoc -I$(GOPATH)/src -Ipkg/cap pkg/cap/cap.proto --go_out=plugins=grpc:pkg/cap
	protoc -I$(GOPATH)/src -Ipkg/alerts pkg/alerts/alerts.proto --go_out=plugins=grpc:pkg/alerts
	protoc -I$(GOPATH)/src -Ipkg/resources pkg/resources/resources.proto --go_out=plugins=grpc:pkg/resources

base: proto
	docker build -t nexus-docker.zacharyseguin.ca/alerts/base:latest .

alerts: base
	docker build -f Dockerfile.alerts -t nexus-docker.zacharyseguin.ca/alerts/alerts:latest  .

resources: base
	docker build -f Dockerfile.resources -t nexus-docker.zacharyseguin.ca/alerts/resources:latest  .

gateway: base
	docker build -f Dockerfile.gateway -t nexus-docker.zacharyseguin.ca/alerts/gateway:latest  .

push: build
	docker push nexus-docker.zacharyseguin.ca/alerts/alerts:latest
	docker push nexus-docker.zacharyseguin.ca/alerts/resources:latest
	docker push nexus-docker.zacharyseguin.ca/alerts/gateway:latest
