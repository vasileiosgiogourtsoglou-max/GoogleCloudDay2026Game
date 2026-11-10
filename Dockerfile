# Build stage: Compile WebAssembly game and server
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 1. Compile Ebitengine game to WebAssembly
RUN mkdir -p /app/public && \
    GOOS=js GOARCH=wasm go build -o /app/public/game.wasm ./cmd/game && \
    cp index.html /app/public/ && \
    if [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
      cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /app/public/wasm_exec.js; \
    else \
      cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" /app/public/wasm_exec.js; \
    fi

# 2. Compile static HTTP server
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 2: Minimal Distroless container for Cloud Run
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /server /app/server
COPY --from=builder /app/public /app/public

EXPOSE 8080
ENTRYPOINT ["/app/server"]
