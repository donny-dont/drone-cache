# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-cache .

FROM alpine:3.4

RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

ADD drone-cache /bin/
ENTRYPOINT ["/bin/drone-cache"]
CMD ["s3"]
