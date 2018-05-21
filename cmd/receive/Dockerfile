FROM alerts:builder AS builder

# Generate final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/receive .
ENTRYPOINT ["./receive"]
