# Copyright 2019 Core Services Team.

FROM --platform=linux/amd64 golang:1.21-alpine as builder

RUN apk add --no-cache ca-certificates git

WORKDIR /eral-promo-library-go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go install .

FROM --platform=linux/amd64 alpine:3.13
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin /bin

RUN mkdir /filestore
RUN chmod -R 0777 /filestore
COPY /filestore /filestore

RUN mkdir /html
COPY /html /html

RUN mkdir /data
ARG version
RUN echo "$version" >> /data/.version

USER nobody:nobody
ENTRYPOINT ["/bin/eral-promo-library-go"]