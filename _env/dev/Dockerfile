FROM golang:1.25-alpine

ARG HOST_UID=1000
ARG HOST_GID=1000

RUN addgroup -g ${HOST_GID} -S appgroup && \
  adduser -u ${HOST_UID} -G appgroup -S -D -h /home/appuser appuser

RUN mkdir -p /app && chown -R appuser:appgroup /app
WORKDIR /app

RUN apk add --no-cache bash curl git unzip && \
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
  go install github.com/pressly/goose/v3/cmd/goose@latest && \
  go install github.com/air-verse/air@latest && \
  go install github.com/smartystreets/goconvey@latest && \
  go install github.com/securego/gosec/v2/cmd/gosec@latest && \
  go install mvdan.cc/gofumpt@latest

ENV PATH=$PATH:/go/bin

RUN chown -R appuser:appgroup /go

USER appuser

COPY --chown=appuser:appgroup go.mod go.sum ./

RUN go mod download

RUN chmod 444 go.mod go.sum

EXPOSE 8080

CMD ["sh", "-c", "mkdir -p tmp/air && air -c .air.toml"]
