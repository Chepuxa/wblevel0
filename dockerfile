FROM golang:1.22.3-alpine

WORKDIR /app
COPY  go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -C cmd -o level0
CMD ["./cmd/level0"]