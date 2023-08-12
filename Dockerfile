# syntax=docker/dockerfile:1

# не используятся stage build-ы — их следует использовать,
# чтобы итоговый образ, который будет дистрибьютится, весил как можно меньше
# версия отстает от последней на 2 мажорных
FROM golang:1.19 as builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

# зачем вам тут goose?
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go mod download
RUN go mod tidy

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy

COPY . .
# выглядит очень неаккуратно, зачем это?
# почему не go build cmd/main.go?
COPY cmd/*.go ./

# Build
# сомнительное имя для бинаря
RUN CGO_ENABLED=0 GOOS=linux go build -o /cmd
#RUN go build -o cmd/main

# Run
CMD ["/cmd"]
