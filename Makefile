# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Docker env
DOCKERFILES := $(shell find . -type f -name '*Dockerfile*' !  -path "./debug/*" )
LINTER_VERSION := v2.12.0

PKG := github.com/outscale/k8s-image-preloader
IMAGE := outscale/k8s-image-preloader
IMAGE_TAG ?= $(shell git describe --tags --always --dirty)
VERSION ?= ${IMAGE_TAG}
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS ?= "-s -w -X ${PKG}/cmd.version=${VERSION} -X ${PKG}/cmd.buildDate=${BUILD_DATE}"
TRIVY_IMAGE := aquasec/trivy:0.30.0


all: help

.PHONY: help
help:
	@echo "help:"
	@echo "  - image-build        : build Docker image"
	@echo "  - image-tag          : tag Docker image"
	@echo "  - image-push         : push Docker image"
	@echo "  - dockerlint         : check Dockerfile"

.PHONY: build
build:
	CGO_ENABLED=0 go build -a -ldflags ${LDFLAGS} -o preloader main.go

.PHONY: image-build
image-build:
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(IMAGE_TAG) .

REGISTRY_IMAGE ?= $(IMAGE)
REGISTRY_TAG ?= $(IMAGE_TAG)
image-tag:
	docker tag $(IMAGE):$(IMAGE_TAG) $(REGISTRY_IMAGE):$(REGISTRY_TAG)

image-push:
	docker push $(REGISTRY_IMAGE):$(REGISTRY_TAG)

.PHONY: dockerlint
dockerlint:
	@echo "Lint images =>  $(DOCKERFILES)"
	$(foreach image,$(DOCKERFILES), echo "Lint  ${image} " ; docker run --rm -i hadolint/hadolint:${LINTER_VERSION} hadolint --info DL3008 -t warning - < ${image} || exit 1 ; )

.PHONY: trivy-scan
trivy-scan:
	docker pull $(TRIVY_IMAGE)
	docker run --rm \
			-v /var/run/docker.sock:/var/run/docker.sock \
			-v ${PWD}/.trivyignore:/root/.trivyignore \
			-v ${PWD}/.trivyscan/:/root/.trivyscan \
			$(TRIVY_IMAGE) \
			image \
			--exit-code 1 \
			--severity="HIGH,CRITICAL" \
			--ignorefile /root/.trivyignore \
			--security-checks vuln \
			--format sarif -o /root/.trivyscan/report.sarif \
			$(IMAGE):$(IMAGE_TAG)
