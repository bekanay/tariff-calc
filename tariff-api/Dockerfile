# --- build stage ---
ARG GO_VERSION=1.24.4
FROM golang:${GO_VERSION}-bookworm AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
WORKDIR /src

# Build dependencies for CGO + IBM DB2 CLI download script
RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential ca-certificates curl unzip tar && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

# Download IBM DB2 CLI driver used by go_ibm_db and place it at a fixed location
RUN mkdir -p /opt/ibm && \
    IBMDB_DIR=$(go env GOPATH)/pkg/mod/github.com/ibmdb/go_ibm_db@$(go list -m -f '{{.Version}}' github.com/ibmdb/go_ibm_db) && \
    cd ${IBMDB_DIR}/installer && go run setup.go && \
    cp -r /go/pkg/mod/github.com/ibmdb/clidriver /opt/ibm/clidriver

# Add clidriver include/lib paths for CGO
ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV CGO_CFLAGS=-I${IBM_DB_HOME}/include
ENV CGO_LDFLAGS=-L${IBM_DB_HOME}/lib

COPY . .
RUN go build -ldflags="-s -w" -o /app/app ./cmd && \
    go build -ldflags="-s -w" -o /app/syncstan ./cmd/sync_stan

# --- runtime stage ---
FROM debian:bookworm-slim
WORKDIR /app

RUN useradd -m appuser && \
    apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates libxml2 libstdc++6 && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy app binary and DB2 CLI driver
COPY --from=build /app/app /usr/local/bin/app
COPY --from=build /app/syncstan /usr/local/bin/syncstan
COPY --from=build /opt/ibm/clidriver /opt/ibm/clidriver

ENV IBM_DB_HOME=/opt/ibm/clidriver
ENV LD_LIBRARY_PATH=${IBM_DB_HOME}/lib
ENV DB2CODEPAGE=1208
ENV LANG=C.UTF-8

USER appuser
EXPOSE 8081
CMD ["app"]
