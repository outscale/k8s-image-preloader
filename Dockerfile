FROM debian:12@sha256:813017f3d62be4b5891a7acca6a01bdcd4b8513daa81b1ab99d3a50385b26931
ARG TARGETPLATFORM
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY $TARGETPLATFORM/preloader /usr/bin/
COPY snapshot.sh /
ENTRYPOINT ["/snapshot.sh"]
