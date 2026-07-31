# syntax=docker/dockerfile:1

# ─────────────────────────────────────────────────────────────
# Build stage: install deps, generate templ code, compile binary
# ─────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS build

# git is occasionally needed for module fetching / build stamping.
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies first (re-runs only when go.mod/go.sum change).
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Install the templ CLI at the version pinned in go.mod. Placed BEFORE the
# source copy so it stays cached across code changes. The @version form
# resolves the CLI's OWN dependencies independently of this project's go.mod —
# plain `go run github.com/a-h/templ/cmd/templ` fails because the CLI pulls in
# modules (LSP, file watcher, etc.) that aren't in the app's go.sum.
# If you bump templ in go.mod, bump this version to match.
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/a-h/templ/cmd/templ@v0.3.1020

# Copy the rest of the source.
COPY . .

# Generate *_templ.go from *.templ, then build a static binary (CGO off).
# The cache mounts keep repeat builds fast — this matters for `docker compose
# watch`, which rebuilds on every save.
RUN templ generate
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# ─────────────────────────────────────────────────────────────
# Runtime stage: tiny image with just the binary + static assets
# ─────────────────────────────────────────────────────────────
FROM alpine:3.20 AS runtime

# Run as a non-root user.
RUN adduser -D -u 10001 appuser

WORKDIR /app

# The compiled server.
COPY --from=build /app/server /app/server

# Static assets (css / js). Served from disk by main.go at /public/,
# so they must be present in the image.
COPY --from=build /app/public /app/public

USER appuser

# Port the server listens on INSIDE the container.
# main.go currently hardcodes ":3000". PORT below only takes effect if you
# update main.go to read it (see README > "Making the port configurable").
ENV PORT=3000
EXPOSE 3000

ENTRYPOINT ["/app/server"]
