FROM docker.io/python:3.12-bullseye AS build
WORKDIR "/app"

# hadolint ignore=DL3008,DL3013
RUN set -eux && \
    apt-get update; \
    apt-get install --no-install-recommends -y \
        python3-dev build-essential cmake libdbus-1-dev dbus-daemon \
        libgirepository1.0-dev libsystemd-dev libgpgme-dev libvirt-dev upx; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*; \
    python3 -m pip install --no-cache-dir --upgrade --force --ignore-installed pip; \
    python3 -m pip install --no-cache-dir --upgrade wheel pyinstaller; \
    python3 -m pip install --no-cache-dir --upgrade asyncio_simple_http_server;

RUN ls /home

ADD data ./
RUN ls /app

COPY requirements.txt ./
#RUN python3 -m pip install bcc==0.29.1 
RUN python3 -m pip install --no-cache-dir --requirement requirements.txt

COPY natasha1.py natasha1.spec ./

RUN pyinstaller natasha1.spec && \
    strip -s -R .comment -R .gnu.version --strip-unneeded dist/natasha1;

# Собираем итоговый образ
FROM docker.io/debian:bullseye-slim

RUN set -eux && \
    apt-get update; \
    apt-get install --no-install-recommends -y \
        busybox curl; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*;

# Copy components
COPY --from=build /app/dist/ /
# COPY config.yml /

ADD data ./

RUN ls /

EXPOSE 8888

ENTRYPOINT ["/natasha1"]
# CMD ["/app/config.yml"]
