FROM golang:1.22-alpine AS build_finala

# Use --no-cache and install only necessary packages
RUN apk add --no-cache git make gcc musl-dev && \
	git config --global http.https://gopkg.in.followRedirects true 

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

RUN make build-linux && \
	mv /app/finala_linux /app/finala


FROM node:18-alpine AS build_ui

# Use --no-cache
RUN apk add --no-cache make

WORKDIR /app

# Copy Makefile first if needed by make build-ui from root
COPY Makefile .

# Copy package.json and package-lock.json to ui directory
COPY ui/package.json ui/package-lock.json* ./ui/

# Install UI dependencies
RUN cd ui && npm install --legacy-peer-deps

# Copy the rest of the UI source code
COPY ui/ ./ui/

# Run the UI build (make sure Makefile handles paths correctly or adjust make command)
RUN make build-ui


# Final stage with updated Alpine
FROM alpine:3.19

# Add non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Create directories and set ownership for appuser
RUN mkdir -p /app/ui/build && \
    chown -R appuser:appgroup /app

# Copy UI assets into /app/ui/build, owned by appuser
COPY --from=build_ui --chown=appuser:appgroup /app/ui/build /app/ui/build

# Copy the finala binary
COPY --from=build_finala /app/finala /bin/finala
RUN chmod +x /bin/finala # Ensure executable

# Set the working directory for the application
WORKDIR /app

# Set the user
USER appuser

ENTRYPOINT ["/bin/finala"]