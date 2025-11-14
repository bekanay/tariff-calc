# --- build stage ---
ARG GO_VERSION=1.24.4
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/app ./cmd

# --- runtime stage ---
FROM alpine:3.20
WORKDIR /app
RUN adduser -D -g '' appuser && apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build /app/app /usr/local/bin/app
USER appuser
EXPOSE 8081
CMD ["app"]
