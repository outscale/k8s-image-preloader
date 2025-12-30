# Build the manager binary
FROM golang:1.25@sha256:31c1e53dfc1cc2d269deec9c83f58729fa3c53dc9a576f6426109d1e319e9e9a AS builder
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

FROM debian:12@sha256:c66c66fac809bfb56a8001b12f08181a49b6db832d2c8ddabe22b6374264055f
WORKDIR /
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY snapshot.sh .
COPY --from=builder /workspace/preloader .
ENTRYPOINT ["/snapshot.sh"]
