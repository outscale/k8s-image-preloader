FROM debian:12@sha256:30482e873082e906a4908c10529180aefb6f77620aea7404b909829fadc5d168
ARG TARGETPLATFORM
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY $TARGETPLATFORM/preloader /usr/bin/
COPY snapshot.sh /
ENTRYPOINT ["/snapshot.sh"]
