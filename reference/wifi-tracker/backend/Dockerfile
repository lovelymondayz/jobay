# ---------- Build Stage ----------
FROM golang:1.24 AS build-stage

# Build arg untuk menentukan entry point (misalnya: cmd/server)
ARG APP_DIR=cmd/server

# Set working directory
WORKDIR /app

# Copy Go modules & download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all source code
COPY . .

# Build Go binary dari direktori yang dipilih
RUN go build -o app ./${APP_DIR} && chmod +x app


# ---------- Production Stage ----------
FROM alpine:latest AS production-stage

# Install packages: C lib compatibility & timezone
RUN apk add --no-cache libc6-compat tzdata

# Set timezone (opsional)
ENV TZ=Asia/Jakarta

# Set working directory
WORKDIR /app

# Copy hasil build binary
COPY --from=build-stage /app/app /app/app

# (Opsional) Copy resource static bila diperlukan
COPY --from=build-stage /app/docs /app/docs
COPY --from=build-stage /app/public /app/public

# Port sesuai yang digunakan di aplikasi
EXPOSE 8081

# Jalankan binary
CMD ["/app/app"]
