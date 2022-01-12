FROM golang:1.16 AS build
WORKDIR /src
COPY . .
RUN GOPRIVATE=git.maynitek.ru CGO_ENABLED=0 go install ./... && \
    CGO_ENABLED=0 go get github.com/go-delve/delve/cmd/dlv

FROM alpine:latest

COPY --from=build /go/bin/sstcloud-alice-gateway /srv/sstcloud-alice-gateway
COPY --from=build /go/bin/dlv /srv/dlv

ENTRYPOINT ["/srv/sstcloud-alice-gateway"]