
ARG BUILDER_IMAGE=default-route-openshift-image-registry.apps.ocp-new-dev.bri.co.id/bricams/golang:1.21.1-alpine
ARG BASE_IMAGE=default-route-openshift-image-registry.apps.ocp-new-dev.bri.co.id/bricams/bitnami-minideb:buster

FROM $BUILDER_IMAGE as builder

ENV http_proxy 'http://proxy4.bri.co.id:1707'
ENV https_proxy 'http://proxy4.bri.co.id:1707'

COPY . /root/go/src/app/

ARG BUILD_VERSION=1.0.0
ARG GOPROXY
ARG GOSUMDB=sum.golang.org

WORKDIR /root/go/src/app

ENV PATH="${PATH}:/usr/local/go/bin"
ENV BUILD_VERSION=$BUILD_VERSION
ENV GOPROXY=$GOPROXY
ENV GOSUMDB=$GOSUMDB

RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -ldflags "-X main.version=$(BUILD_VERSION)" -installsuffix cgo -o app ./server

ENV http_proxy=
ENV https_proxy=

FROM $BASE_IMAGE

WORKDIR /usr/app

COPY --from=builder /root/go/src/app/app /usr/app/app
COPY --from=builder /root/go/src/app/www /usr/app/www
COPY --from=builder /root/go/src/app/grpc_health_probe-linux-amd64 /usr/app/grpc_health_probe-linux-amd64
RUN chmod a+x /usr/app/grpc_health_probe-linux-amd64

LABEL authors="Qcash Squad Name"

# PotatoBeans Co. adheres to OCI image specification.
LABEL org.opencontainers.image.authors="Qcash Squad Name"
LABEL org.opencontainers.image.title="go-base"
LABEL org.opencontainers.image.description="Qcash GRPC Feature Name"
LABEL org.opencontainers.image.vendor=""

EXPOSE 9090
EXPOSE 3000

ENTRYPOINT ["/usr/app/app"]
CMD ["grpc-gw-server", "--port1", "9090", "--port2", "3000", "--grpc-endpoint", ":9090"]
