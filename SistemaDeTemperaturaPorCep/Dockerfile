FROM golang:1.23.1-alpine3.19 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/server

FROM busybox:musl AS mybusybox

FROM debian:stable-slim AS mydebian
RUN apt update && apt install -y ca-certificates

FROM scratch AS base

COPY --from=mybusybox /etc/passwd /etc/passwd
COPY --from=mybusybox /etc/group /etc/group
COPY --from=mydebian /etc/ssl/certs/ca-certificates.crt /app/ssl/certs/

RUN --mount=from=mybusybox,dst=/usr/ ["busybox", "sh", "-c", "mkdir -p /app && chmod 777 /app"]
RUN --mount=from=mybusybox,dst=/usr/ ["busybox", "sh", "-c", "addgroup -S go -g 1000 && adduser -S go -u 1000 --ingroup go --disabled-password"]

ENV HOME=/app
ENV USER=go
ENV PATH=/usr/local/bin:/app
ENV SSL_CERT_DIR=/app/ssl/certs

FROM base AS release

WORKDIR /app
USER go

COPY --from=builder /app/myapp /app/myapp
COPY --from=debian:stable-slim /lib/x86_64-linux-gnu /lib/x86_64-linux-gnu
COPY --from=debian:stable-slim /lib64/ld-linux-x86-64.so.2 /lib64/ld-linux-x86-64.so.2
COPY --from=mydebian /bin/sh /app/sh
COPY --from=mydebian /usr/bin/sleep /app/sleep
COPY cmd/server/.env /app/.env
COPY container/app-run.sh /app/app-run.sh

CMD [ "app-run.sh" ]