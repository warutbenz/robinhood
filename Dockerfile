FROM golang:alpine

COPY . /app
WORKDIR /app
RUN go install github.com/cespare/reflex@latest

EXPOSE 8080
