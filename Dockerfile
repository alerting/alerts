FROM golang:1.10-alpine

WORKDIR /go/src/github.com/alerting/alerts

# Install required software
RUN apk --update add git
RUN go get github.com/golang/dep/cmd/dep

# Install dependencies
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

# Copy source code
COPY ./ ./

# Run the build
RUN CGO_ENABLED=0 GOOS=linux go install ./cmd/...