# syntax=docker/dockerfile:1

FROM golang:1.19 as build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

#RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod download
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

COPY . .
COPY cmd/*.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /cmd
#RUN go build -o cmd/main

# Run
CMD ["/cmd"]

FROM alpine:latest as production
COPY --from=build cmd/ .
COPY .env .
RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

CMD ["/cmd"]
