FROM --platform=${TARGETPLATFORM:-linux/amd64} ghcr.io/openfaas/of-watchdog:0.8.4 as watchdog
FROM golang:1.16.3-alpine3.13 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG GO111MODULE="on"
ENV CGO_ENABLED=0

RUN apk --no-cache add git

COPY --from=watchdog /fwatchdog /usr/bin/fwatchdog
RUN chmod +x /usr/bin/fwatchdog

RUN mkdir -p /go/src/handler
WORKDIR /go/src/handler
COPY . .

# Add user overrides to the root go.mod, which is the only place "replace" can be used
RUN cat function/GO_REPLACE.txt >> ./go.mod || exit 0

RUN go mod tidy 

# Run go test, gofmt, go vet
RUN chmod +x test.sh \
    && ./test.sh \
    && go build --ldflags "-s -w" -a -installsuffix cgo -o handler .

FROM alpine:latest
# Add non root user and certs
RUN apk --no-cache add ca-certificates \
    && addgroup -S app && adduser -S -g app app
# Split instructions so that buildkit can run & cache 
# the previous command ahead of time.
RUN mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

COPY --from=build /go/src/handler/handler /usr/bin/fwatchdog /go/src/handler/function/ ./
RUN chown -R app /home/app

USER app
ENV fprocess="./handler" mode="http" upstream_url="http://127.0.0.1:8082"
CMD ["./fwatchdog"]
