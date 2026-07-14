FROM debian:12@sha256:9344f8b8992482f80cba753f323adeaf17690076c095ccff6cc9536be98185dc
ARG TARGETPLATFORM
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY $TARGETPLATFORM/preloader /usr/bin/
COPY snapshot.sh /
ENTRYPOINT ["/snapshot.sh"]
