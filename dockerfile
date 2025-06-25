# ────────────────────────────────
# 1️⃣  build stage – compiles loghub
# ────────────────────────────────
FROM golang:1.24.4-alpine AS build

# Install build tools (optional: git for version info)
RUN apk add --no-cache git

# Enable reproducible builds
ENV CGO_ENABLED=0 \
    GOFLAGS="-trimpath" \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Compile only the CLI (main.go lives in repo root)
RUN go build -o /loghub .

# ────────────────────────────────
# 2️⃣  runtime stage – tiny image
# ────────────────────────────────
FROM alpine:3.19

# Add a non-root user for safety
RUN addgroup -S app && adduser -S app -G app
USER app

# Copy the statically-linked binary from build stage
COPY --from=build /loghub /usr/local/bin/loghub

# Default working directory for relative paths
WORKDIR /app

# Create a writable volume for logs (host can mount over it)
VOLUME ["/logs"]

# By default just print the help text;
# override CMD in `docker run` to use watch/filter/export.
CMD ["loghub", "--help"]

CMD tail -f /dev/null