
# Note: You need an individual entry below for them to be built
COMMANDS=$(wildcard cmd/*)
TAG := latest

build: $(COMMANDS)

.PHONY: build clean builder $(COMMANDS)

clean:
	-docker rmi alerts:builder alerts:alert

builder:
	docker build -t alerts:builder .

cmd/alert: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-${TAG} $@

cmd/gateway: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-${TAG} $@

cmd/import: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-${TAG} $@

cmd/receive: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-${TAG} $@
