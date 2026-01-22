FROM debian:12@sha256:23f8d499dcf5559a92ef814c309d747268420cf5e5e60f067f495d509f0aeb93

RUN apt-get update && apt-get upgrade -y && apt-get install -y curl gpg

RUN curl -fsSL https://packages.buildkite.com/helm-linux/helm-debian/gpgkey | gpg --dearmor > /usr/share/keyrings/helm.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" > /etc/apt/sources.list.d/helm-stable-debian.list

RUN apt-get update && apt-get install -y helm && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -LO https://dl.k8s.io/release/v1.30.12/bin/linux/amd64/kubectl && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

COPY main.sh cleanup.sh snapshot.yaml /

ENTRYPOINT ["/main.sh"]