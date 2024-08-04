FROM golang:1.22-alpine
WORKDIR /app

COPY *.go go.mod go.sum ./
RUN go mod download
RUN go build -o main .

EXPOSE 7000
CMD ["./main"]
