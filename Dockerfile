# Use Go 1.23 alpine image
FROM golang:1.23-alpine

LABEL maintainer="Teodor Savin <teodorsavin@gmail.com>"

# Install necessary dependencies including Air
RUN apk update && apk add --no-cache git ca-certificates && go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source files
COPY . .

EXPOSE 8080

# Run Air for hot-reloading
CMD ["air"]
