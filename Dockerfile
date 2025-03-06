FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY srvctl /usr/local/bin/srvctl

ENTRYPOINT ["/usr/local/bin/srvctl"]
