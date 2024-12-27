# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o dnd ./cmd/server

# Stage 2: Create a lightweight image for running
FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /app

# Copia o binário
COPY --from=builder /app/dnd .

# (NOVO) Copia também o config.yaml para /app
COPY config.yaml .

# Copia wait-for-it.sh
COPY /scripts/wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./dnd"]
