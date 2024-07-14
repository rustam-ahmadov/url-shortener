FROM golang:1.22-alpine3.19
WORKDIR /app
COPY . .
RUN go build -o url-shortener cmd/url-shortener/main.go
EXPOSE 8080
CMD ["./url-shortener"]