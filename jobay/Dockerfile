# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o jobay .

# Frontend build
FROM node:20-slim AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Final
FROM alpine:3.19
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy Go binary
COPY --from=builder /app/jobay .

# Copy frontend
COPY --from=frontend /app/frontend/dist ./frontend/dist

# Data directory
RUN mkdir -p /app/data/cvs

# Environment
ENV PORT=8080
ENV DB_PATH=/app/data/jobay.json
ENV UPLOADS_DIR=/app/data/cvs

EXPOSE 8080

CMD ["./jobay"]
