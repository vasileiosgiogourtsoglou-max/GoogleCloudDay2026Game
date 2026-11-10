# Build stage: Compile WebAssembly game and embedded HTTP server
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 1. Compile Ebitengine game to WebAssembly
RUN mkdir -p public && \
    GOOS=js GOARCH=wasm go build -o public/game.wasm ./cmd/game && \
    cp index.html public/ && \
    if [ -f "$(go env GOROOT)/lib/wasm/wasm_exec.js" ]; then \
      cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" public/wasm_exec.js; \
    else \
      cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" public/wasm_exec.js; \
    fi

# 2. Compile self-contained HTTP server (embeds public/*)
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Stage 2: Minimal Distroless container for Cloud Run
FROM gcr.io/distroless/static-debian12

COPY --from=builder /server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
