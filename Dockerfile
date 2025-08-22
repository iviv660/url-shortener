# ---------- Build stage ----------
FROM golang:1.24.4-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем пакет main из каталога ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -trimpath -ldflags "-s -w" -o /app ./cmd 

# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /app /app
ENV GIN_MODE=release
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app"]

