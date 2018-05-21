
# Note: You need an individual entry below for them to be built
COMMANDS=$(wildcard cmd/*)

build: $(COMMANDS)

.PHONY: build clean builder $(COMMANDS)

clean:
	-docker rmi alerts:builder alerts:alert

builder:
	docker build -t alerts:builder .

cmd/alert: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-latest $@

cmd/gateway: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-latest $@

cmd/import: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-latest $@

cmd/receive: builder
	docker build -t zachomedia/alerting:$(subst cmd/,,$@)-latest $@
