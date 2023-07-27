# syntax=docker/dockerfile:1

FROM golang:1.19 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod download
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

COPY . .
COPY cmd/*.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /cmd
#RUN go build -o cmd/main

# Run
CMD ["/cmd"]
