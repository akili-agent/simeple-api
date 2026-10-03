# syntax=docker/dockerfile:1

########## Build stage ##########
FROM golang:1.26-alpine AS build

WORKDIR /src

# Download dependencies first so the module layer is cached independently of
# the source code.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary: CGO disabled so the runtime image needs no libc.
RUN CGO_ENABLED=0 GOOS=linux go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/server ./cmd/server

########## Runtime stage ##########
FROM alpine:3.21

# CA certificates for outbound TLS, and a dedicated non-root user.
RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 app

COPY --from=build /out/server /usr/local/bin/server

USER app
EXPOSE 8080

ENV PORT=8080 \
    LOG_LEVEL=info

ENTRYPOINT ["server"]
