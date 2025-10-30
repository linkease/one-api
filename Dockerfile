FROM --platform=$BUILDPLATFORM node:16 AS builder

WORKDIR /web
COPY ./VERSION .
COPY ./web .

RUN npm install --prefix /web/default & \
    npm install --prefix /web/berry & \
    npm install --prefix /web/air & \
    wait

RUN DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ./VERSION) npm run build --prefix /web/default && \
    DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ./VERSION) npm run build --prefix /web/berry && \
    DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat ./VERSION) npm run build --prefix /web/air

# Debug: verify front-end build artifacts are where we expect them
RUN echo "=== Verify web/build artifacts ===" && \
        echo "PWD: $(pwd)" && \
        echo "--- /web ---" && ls -la /web || true && \
        echo "--- /web/build ---" && ls -la /web/build || true && \
        for T in default berry air; do \
            echo "--- Theme: $T ---"; \
            if [ -d "/web/build/$T" ]; then \
                ls -la "/web/build/$T"; \
                if [ -f "/web/build/$T/index.html" ]; then \
                    echo "OK: /web/build/$T/index.html exists"; \
                    # Print first few referenced assets for quick sanity check (non-fatal)
                    head -n 60 "/web/build/$T/index.html" | grep -Eo 'src="/kc-admin/[^" ]+\.js"|href="/kc-admin/[^" ]+\.css"' | head -n 5 || true; \
                else \
                    echo "MISSING: /web/build/$T/index.html"; \
                fi; \
            else \
                echo "MISSING DIR: /web/build/$T"; \
            fi; \
        done && \
        echo "--- index.html discovered under /web/build (maxdepth=2) ---" && \
        find /web/build -maxdepth 2 -type f -name 'index.html' -print || true

FROM golang:alpine AS builder2

RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev \
    build-base

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=builder /web/build ./web/build

RUN go build -trimpath -ldflags "-s -w -X 'github.com/songquanpeng/one-api/common.Version=$(cat VERSION)' -linkmode external -extldflags '-static'" -o one-api

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder2 /build/one-api /

EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/one-api"]