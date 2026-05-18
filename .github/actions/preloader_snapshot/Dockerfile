FROM debian:12@sha256:85019db29298555fd1a5f4bb57673ae989414a9884117c75d7a3e1a6cce21688

RUN apt-get update && apt-get upgrade -y && apt-get install -y curl gpg

RUN curl -fsSL --retry 5 --retry-all-errors https://packages.buildkite.com/helm-linux/helm-debian/gpgkey | gpg --dearmor > /usr/share/keyrings/helm.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" > /etc/apt/sources.list.d/helm-stable-debian.list

RUN apt-get update && apt-get install -y helm=3.20.0-1 && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -sSLO --retry 5 --retry-all-errors https://dl.k8s.io/release/v1.34.8/bin/linux/amd64/kubectl && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

COPY main.sh cleanup.sh snapshot.yaml /

ENTRYPOINT ["/main.sh"]
