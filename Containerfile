## Build image
FROM golang:alpine AS build

WORKDIR /app

COPY . /app

RUN go mod download
RUN go build -o /astartebot

## Deployable image

FROM alpine

WORKDIR /
COPY --from=build /astartebot /astartebot
EXPOSE 8080
USER nobody:nobody

ENTRYPOINT ["/astartebot"]
