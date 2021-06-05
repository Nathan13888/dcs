# IMAGE: BUILDER
FROM golang:1.16-alpine as builder
WORKDIR /build/
COPY . .
RUN apk add make
RUN apk add gcc musl-dev
RUN make build
#RUN go build -o ./bin/dcs

# IMAGE: CONTAINER
FROM alpine:latest
MAINTAINER Nathan13888
WORKDIR /app

# ENV VARIABLES
#ENV PORT=6969
#ENV TZ=America/Toronto
ENV DOWNLOADPATH=/app/downloads
ENV USER=dcs
ENV UID=12345
ENV GID=23456
## TIMEZONE
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/America/Toronto /etc/localtime
## CREATE GROUP
RUN addgroup --gid $GID $USER
## CREATE USER
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

RUN mkdir -p .dcs/logs downloads
RUN chown -R dcs:dcs .
RUN chmod -R 755 .
USER dcs

COPY --from=builder --chown=dcs /build/bin/dcs /app/dcs

VOLUME /app/.dcs
VOLUME /app/downloads
EXPOSE 6969

ENTRYPOINT [ "/app/dcs", "service", "start" ]

