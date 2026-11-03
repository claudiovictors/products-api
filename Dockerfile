# ---- Build stage ---------------------------------------------------------
# Use a Go image that matches the version declared in go.mod to compile the app.
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy the module definition first so the dependency download layer is cached
# and only re-runs when go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build a static linux binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api .

# ---- Runtime stage -------------------------------------------------------
# The final image only contains the compiled binary, keeping it small.
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
