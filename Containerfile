##
## -- runtime environment
##

FROM    golang:1.24-alpine3.21 AS build-env

#       https://github.com/docker-library/official-images#multiple-architectures
#       https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG     TARGETPLATFORM
ARG     TARGETOS
ARG     TARGETARCH

ARG     VERSION=latest

ADD     . /src/speedtest-exporter
RUN     apk add -U --no-cache make git
RUN     make LDFLAGS="-ldflags -w" -C /src/speedtest-exporter bin/speedtest-exporter-$VERSION.$TARGETOS.$TARGETARCH
RUN     go install -v github.com/showwin/speedtest-go@latest

##
## -- runtime environment
##

FROM    alpine:3.21 AS rt-env

RUN     apk add -U --no-cache tini ca-certificates && apk del apk-tools libc-utils
COPY    --from=build-env /src/speedtest-exporter/bin/* /usr/bin/speedtest-exporter
COPY    --from=build-env /go/bin/speedtest-go /usr/bin/speedtest-go

EXPOSE  8080
ENTRYPOINT ["/sbin/tini", "--", "/usr/bin/speedtest-exporter"]
