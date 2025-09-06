# ---- Builder ----
FROM golang:1.21-alpine AS build
# Auto-fetch the right Go toolchain for your 'go' directive
ENV GOTOOLCHAIN=auto CGO_ENABLED=0

WORKDIR /app

# Cache module download separately
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest and build with build cache
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o server ./cmd/server

# ---- Runtime ----
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]
