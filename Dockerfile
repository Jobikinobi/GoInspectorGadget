#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
# Build the investigator CLI application
RUN go build -o /go/bin/app -v ./cmd/investigator

#final stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates curl

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /go/bin/app /app/
RUN chmod +x /app/app

# Create data directory
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Simple healthcheck using the application's help command
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD /app/app help > /dev/null || exit 1

# Set entrypoint with help as default command
ENTRYPOINT ["/app/app"]
CMD ["help"]

LABEL Name=goinspectorgadget Version=0.0.1
EXPOSE 3000
