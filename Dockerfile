ARG DOCKER_PROXY

FROM golang:1.17-alpine3.15 AS builder

ARG GOSUMDB
ARG GOPRIVATE
ARG GOPROXY

ENV APPDIR ${GOPATH}/src/api
ENV ARTIFACT /build/api

RUN apk --update --no-cache add tzdata git

RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" >  /etc/timezone
RUN date
RUN apk del tzdata

RUN mkdir -p ${APPDIR}
WORKDIR ${APPDIR}
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" -o ${ARTIFACT} main.go

FROM alpine:3.15

ENV ARTIFACT /build
ENV BINDIR /usr/local/bin

RUN apk --update --no-cache add tzdata git
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime
RUN echo "Europe/Moscow" >  /etc/timezone
RUN date

COPY --from=builder ${ARTIFACT} ${BINDIR}

WORKDIR ${BINDIR}

RUN mkdir data

ENTRYPOINT [ "api" ]
