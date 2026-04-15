# syntax=docker/dockerfile:1

# ── build stage ──────────────────────────────────────────────────────────────
FROM golang:1.26.2 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o reduct .

# ── runtime stage (scratch) ───────────────────────────────────────────────────
FROM scratch

COPY --from=builder /build/reduct /reduct

ENTRYPOINT ["/reduct"]
