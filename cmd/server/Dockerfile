FROM golang:1.17-alpine3.16 AS build

WORKDIR /pow-tcp-server
COPY . .
RUN go mod download

RUN go build -o /app cmd/server/main.go

FROM alpine

WORKDIR /

COPY --from=build /app /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]