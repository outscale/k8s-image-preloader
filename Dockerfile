FROM debian:12@sha256:854dd8ab51d62866dbd8f5431d15ce77b12983a2a4a1e1fe43a2395dc7ef9f61
ARG TARGETPLATFORM
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY $TARGETPLATFORM/preloader /usr/bin/
COPY snapshot.sh /
ENTRYPOINT ["/snapshot.sh"]
