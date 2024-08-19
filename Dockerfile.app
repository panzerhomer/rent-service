FROM golang:1.22
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main cmd/api/main.go
CMD ["/app/main"] 