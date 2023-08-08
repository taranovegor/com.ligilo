FROM golang:1.21 AS build

WORKDIR /app

COPY ./ /app

RUN go mod download

ARG SCOPE

RUN CGO_ENABLED=0 GOOS=linux go build -v -o $SCOPE ./cmd/$SCOPE

FROM debian:12-slim

WORKDIR /ligilo

ARG SCOPE
ENV SCOPE=$SCOPE

COPY --from=build /app/$SCOPE $SCOPE

ENTRYPOINT /ligilo/$SCOPE
