FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc g++ make git

WORKDIR /go/src/app

ADD . .

RUN go get

RUN go build -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

COPY --chown=65534:65534 --from=builder /go/src/app/app .

USER 65534
