# Multi-stage build for Okta Dev Account Migration Tool (envsync) with proper dependencies
FROM golang:1.24-alpine AS builder

# Install git and make for building dependencies
RUN apk add --no-cache git make openssh-client

# First, build the custom okta-cli-client from the specific fork/branch as per README
WORKDIR /tmp/okta-cli-client
RUN git clone https://github.com/edunham/okta-cli-client.git . && \
    git checkout file-handling && \
    make install

# Now build envsync - copy the local source
WORKDIR /app
COPY . .

# Download Go modules and build envsync using make install as per README
RUN go mod download && \
    make install

# Final stage - lightweight runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests to Okta APIs
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 envsync && \
    adduser -D -u 1001 -G envsync envsync

# Create necessary directories
RUN mkdir -p /home/envsync/.okta && \
    chown -R envsync:envsync /home/envsync

# Copy both binaries from builder stage
COPY --from=builder /go/bin/envsync /usr/local/bin/envsync
COPY --from=builder /go/bin/okta-cli-client /usr/local/bin/okta-cli-client

# Switch to non-root user
USER envsync
WORKDIR /home/envsync

# Set default environment variables
ENV OKTA_CLI_CLIENT_ID=""
ENV OKTA_CLI_ORGURL=""
ENV OKTA_CLI_TOKEN=""

# Create volume for Okta configuration and backup data
VOLUME ["/home/envsync/.okta"]

# Default command shows help
CMD ["envsync", "-h"]
