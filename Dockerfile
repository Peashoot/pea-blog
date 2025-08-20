# Multi-stage build
# Stage 1: Build frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production

COPY frontend/ .
RUN npm run build-only

# Stage 2: Build backend
FROM golang:1.23-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app/backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/server/main.go

# Stage 3: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/main ./
COPY --from=backend-builder /app/backend/.env ./

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Create database directory
RUN mkdir -p /app/data

EXPOSE 8080

ENV DATABASE_URL=/app/data/pea_blog.db

CMD ["./main"]