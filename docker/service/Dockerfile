FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /usr/src/app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build ./cmd/fleetState/main.go

# Copy binary from build to main folder
##RUN cp ./cmd/faceit/faceit .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["./main"]
