# Build the manager binary
FROM golang:1.25@sha256:0f406d34b7cb7255d0700af02ec28a2c88f1e00701055f4c282aa4c3ec0b3245 AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION

WORKDIR /workspace

COPY Makefile Makefile

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} make build

FROM debian:12@sha256:23f8d499dcf5559a92ef814c309d747268420cf5e5e60f067f495d509f0aeb93
WORKDIR /
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY snapshot.sh .
COPY --from=builder /workspace/preloader .
ENTRYPOINT ["/snapshot.sh"]
