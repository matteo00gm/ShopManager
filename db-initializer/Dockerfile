FROM golang

WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o db-initializer .

# Command to run the binary
CMD ["./db-initializer"]