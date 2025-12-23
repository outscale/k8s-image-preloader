FROM debian:12@sha256:7e5bc0e499a8d50cb1e32287944a90b9ec8fd7d500673e75daff3f52882f5798

RUN apt-get update && apt-get upgrade -y && apt-get install -y curl gpg

RUN curl -fsSL https://packages.buildkite.com/helm-linux/helm-debian/gpgkey | gpg --dearmor > /usr/share/keyrings/helm.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" > /etc/apt/sources.list.d/helm-stable-debian.list

RUN apt-get update && apt-get install -y helm && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -LO https://dl.k8s.io/release/v1.30.12/bin/linux/amd64/kubectl && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

COPY main.sh cleanup.sh snapshot.yaml /

ENTRYPOINT ["/main.sh"]