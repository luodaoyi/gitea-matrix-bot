FROM ubuntu:focal-20221130

ENV TZ="Asia/Shanghai"

ARG TARGETOS
ARG TARGETARCH

COPY ./entrypoint.sh /entrypoint.sh

RUN export DEBIAN_FRONTEND="noninteractive" && \
    apt update && apt install -y ca-certificates tzdata && \
    update-ca-certificates && \
    ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure tzdata && \
    chmod +x /entrypoint.sh

WORKDIR /app
COPY dist/app-${TARGETOS}-${TARGETARCH} ./app
COPY config.ini.example /app/config.ini

VOLUME ["/data"]
EXPOSE 9000
ENTRYPOINT ["/entrypoint.sh"]