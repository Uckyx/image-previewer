FROM golang:1.18.1-alpine
RUN apk add --no-cache tzdata
ENV TZ Europe/Moscow
RUN echo -e "https://nl.alpinelinux.org/alpine/v3.5/main\nhttps://nl.alpinelinux.org/alpine/v3.5/community" > /etc/apk/repositories

ENV USER=appuser
ENV UID=10001
# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
