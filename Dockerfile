FROM golang:1.24.1 AS builder

COPY ./src /app/src

WORKDIR /app/src

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vnc-api


FROM gcr.io/distroless/static

COPY --from=builder /app/src/adapters/api/config/authorization /adapters/api/config/authorization
COPY --from=builder /app/src/core/services/resources /core/services/resources
COPY --from=builder /app/src/vnc-api /vnc-api

ENTRYPOINT ["/vnc-api"]
