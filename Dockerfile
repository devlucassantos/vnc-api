FROM golang:1.24.1 AS builder

ARG GITHUB_USERNAME=$GITHUB_USERNAME
ARG GITHUB_ACCESS_TOKEN=$GITHUB_ACCESS_TOKEN

RUN echo "machine github.com\n\tlogin $GITHUB_USERNAME\n\tpassword $GITHUB_ACCESS_TOKEN" >> ~/.netrc

COPY ./src /app/src

WORKDIR /app/src

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vnc-api


FROM gcr.io/distroless/static

COPY --from=builder /app/src/api/config/authorization /api/config/authorization
COPY --from=builder /app/src/core/services/resources /core/services/resources
COPY --from=builder /app/src/vnc-api /vnc-api

ENTRYPOINT ["/vnc-api"]
