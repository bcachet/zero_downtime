# Build stage
FROM docker.io/library/golang:1.25 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and application source code
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY helloworld .
COPY client .
COPY server .

ARG DIRECTORY=server
# Build the gRPC server application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ${DIRECTORY}/main.go

# Final stage
FROM scratch

# Copy the built executable from the builder stage
COPY --from=builder /app/app /app/app

# Set the working directory inside the container
WORKDIR /app

# Command to run the application
ENTRYPOINT [ "/app/app" ]
CMD [ "-delay=0", "-port=50051" ]
